package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"mymall/pkg/cache"
	"mymall/services/order-service/internal/client/merchantrpc"
	"mymall/services/order-service/internal/client/userrpc"
	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type OrderLogic struct {
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{svcCtx: svcCtx}
}

func orderPayAmount(o *model.Order) float64 {
	if o.PayAmount > 0 {
		return o.PayAmount
	}
	return o.TotalAmount
}

func (l *OrderLogic) CouponPreview(ctx context.Context, userID uint64, items []types.CreateOrderItem, userCouponID uint64) (*merchantrpc.MatchResp, error) {
	matchItems, shopID, _, err := l.buildMatchItems(ctx, items)
	if err != nil {
		return nil, err
	}
	if l.svcCtx.MerchantRPC == nil {
		return nil, errors.New("优惠券服务不可用")
	}
	return l.svcCtx.MerchantRPC.MatchCoupons(ctx, userID, shopID, userCouponID, matchItems)
}

func (l *OrderLogic) buildMatchItems(ctx context.Context, items []types.CreateOrderItem) ([]merchantrpc.MatchItem, uint64, float64, error) {
	if len(items) == 0 {
		return nil, 0, 0, errors.New("订单商品不能为空")
	}
	ids := make([]uint64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ProductID)
	}
	resp, err := l.svcCtx.CatalogRPC.BatchGetProducts(ctx, ids)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("查询商品失败: %w", err)
	}
	productByID := map[uint64]struct {
		shopID uint64
		price  float64
		name   string
	}{}
	for _, p := range resp.GetProducts() {
		productByID[p.GetId()] = struct {
			shopID uint64
			price  float64
			name   string
		}{p.GetShopId(), p.GetSalePrice(), p.GetName()}
	}
	var shopID uint64
	var goods float64
	out := make([]merchantrpc.MatchItem, 0, len(items))
	for _, it := range items {
		p, ok := productByID[it.ProductID]
		if !ok {
			return nil, 0, 0, errors.New("部分商品不存在或已下架")
		}
		if shopID == 0 {
			shopID = p.shopID
		} else if p.shopID != shopID {
			return nil, 0, 0, errors.New("不支持跨店下单，请分开结算")
		}
		qty := it.Quantity
		if qty <= 0 {
			return nil, 0, 0, errors.New("购买数量无效")
		}
		price := p.price
		if it.SeckillEntryID > 0 && l.svcCtx.MerchantRPC != nil {
			// 预览用商品价即可；实际下单再校验秒杀
		}
		amt := price * float64(qty)
		goods += amt
		out = append(out, merchantrpc.MatchItem{
			ProductID: it.ProductID, Amount: amt, SeckillEntryID: it.SeckillEntryID,
		})
	}
	return out, shopID, goods, nil
}

func (l *OrderLogic) CreateOrder(ctx context.Context, userID uint64, addressID uint64, items []types.CreateOrderItem, userCouponID uint64) (*model.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}
	if addressID == 0 {
		return nil, errors.New("请选择收货地址")
	}
	if l.svcCtx.Redis == nil {
		return nil, errors.New("库存服务不可用，请稍后重试")
	}
	if l.svcCtx.UserRPC == nil {
		return nil, errors.New("用户服务不可用")
	}
	addr, err := l.svcCtx.UserRPC.GetAddress(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}
	if addr.ReceiverName == "" || addr.ReceiverPhone == "" || addr.FullAddress() == "" {
		return nil, errors.New("收货地址不完整")
	}

	if l.svcCtx.UserRPC != nil {
		u, err := l.svcCtx.UserRPC.GetUser(ctx, userID)
		if err != nil || u.GetStatus() != 1 {
			return nil, errors.New("用户无效")
		}
	}

	ids := make([]uint64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ProductID)
	}

	resp, err := l.svcCtx.CatalogRPC.BatchGetProducts(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}
	productByID := make(map[uint64]*struct {
		id, shopID uint64
		name       string
		price      float64
		defSku     uint64
	}, len(resp.GetProducts()))
	for _, p := range resp.GetProducts() {
		productByID[p.GetId()] = &struct {
			id, shopID uint64
			name       string
			price      float64
			defSku     uint64
		}{p.GetId(), p.GetShopId(), p.GetName(), p.GetSalePrice(), p.GetDefaultSkuId()}
	}

	var total float64
	var shopID uint64
	orderItems := make([]model.OrderItem, 0, len(items))
	stockItems := make([]model.StockItem, 0, len(items))
	redisItems := make([]cache.StockItem, 0, len(items))
	type seckillHold struct {
		entryID uint64
		qty     int
	}
	var seckillHolds []seckillHold
	for _, it := range items {
		p, ok := productByID[it.ProductID]
		if !ok {
			return nil, errors.New("部分商品不存在或已下架")
		}
		if shopID == 0 {
			shopID = p.shopID
		} else if p.shopID != shopID {
			return nil, errors.New("不支持跨店下单，请分开结算")
		}
		skuID := it.SkuID
		if skuID == 0 {
			skuID = p.defSku
		}
		if skuID == 0 {
			return nil, fmt.Errorf("商品 %s 缺少 SKU，无法扣库存", p.name)
		}
		qty := it.Quantity
		if qty <= 0 {
			return nil, errors.New("购买数量无效")
		}
		snap := it.SkuSnapshot
		if snap == "" {
			snap = fmt.Sprintf(`{"sku_id":%d}`, skuID)
		}
		price := p.price
		var seckillEntryID uint64
		if it.SeckillEntryID > 0 {
			if l.svcCtx.MerchantRPC == nil {
				return nil, errors.New("秒杀服务不可用")
			}
			cr, err := l.svcCtx.MerchantRPC.Consume(ctx, it.SeckillEntryID, it.ProductID, qty)
			if err != nil {
				for _, h := range seckillHolds {
					_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
				}
				return nil, err
			}
			price = cr.SeckillPrice
			seckillEntryID = it.SeckillEntryID
			seckillHolds = append(seckillHolds, seckillHold{entryID: it.SeckillEntryID, qty: qty})
		}
		lineAmount := price * float64(qty)
		total += lineAmount
		orderItems = append(orderItems, model.OrderItem{
			ProductID:      p.id,
			SkuID:          skuID,
			ProductName:    p.name,
			SkuSnapshot:    snap,
			Price:          price,
			Quantity:       qty,
			SeckillEntryID: seckillEntryID,
		})
		stockItems = append(stockItems, model.StockItem{ProductID: p.id, SkuID: skuID, Quantity: qty})
		redisItems = append(redisItems, cache.StockItem{SkuID: skuID, Quantity: qty})
	}

	goodsAmount := total
	discountAmount := 0.0
	payAmount := total
	if userCouponID > 0 {
		if l.svcCtx.MerchantRPC == nil {
			for _, h := range seckillHolds {
				if l.svcCtx.MerchantRPC != nil {
					_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
				}
			}
			return nil, errors.New("优惠券服务不可用")
		}
		matchItems := make([]merchantrpc.MatchItem, 0, len(orderItems))
		for _, oi := range orderItems {
			matchItems = append(matchItems, merchantrpc.MatchItem{
				ProductID: oi.ProductID, Amount: oi.Price * float64(oi.Quantity), SeckillEntryID: oi.SeckillEntryID,
			})
		}
		mr, err := l.svcCtx.MerchantRPC.MatchCoupons(ctx, userID, shopID, userCouponID, matchItems)
		if err != nil {
			for _, h := range seckillHolds {
				_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
			}
			return nil, err
		}
		ok := false
		for _, a := range mr.Available {
			if a.UserCouponID == userCouponID && a.Usable {
				discountAmount = a.DiscountAmount
				ok = true
				break
			}
		}
		if !ok {
			for _, h := range seckillHolds {
				_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
			}
			return nil, errors.New("优惠券不可用")
		}
		payAmount = goodsAmount - discountAmount
		if payAmount < 0.01 {
			payAmount = 0.01
			discountAmount = goodsAmount - payAmount
		}
	}

	orderNo := fmt.Sprintf("ORD%s", uuid.NewString()[:8]+time.Now().Format("150405"))
	order := &model.Order{
		OrderNo:         orderNo,
		UserID:          userID,
		ShopID:          shopID,
		TotalAmount:     payAmount,
		GoodsAmount:     goodsAmount,
		DiscountAmount:  discountAmount,
		PayAmount:       payAmount,
		UserCouponID:    userCouponID,
		ReceiverName:    addr.ReceiverName,
		ReceiverPhone:   addr.ReceiverPhone,
		ReceiverAddress: addr.FullAddress(),
		Status:          model.OrderStatusPending,
	}
	if err := l.svcCtx.Repo.Create(ctx, order, orderItems); err != nil {
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
		}
		return nil, err
	}

	couponLocked := false
	if userCouponID > 0 {
		if err := l.svcCtx.MerchantRPC.LockCoupon(ctx, userCouponID, userID, order.ID, discountAmount); err != nil {
			_ = l.svcCtx.Repo.UpdateStatus(ctx, orderNo, model.OrderStatusFailed)
			for _, h := range seckillHolds {
				_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
			}
			return nil, err
		}
		couponLocked = true
	}
	defer func() {
		if couponLocked {
			_ = l.svcCtx.MerchantRPC.UnlockCoupon(context.Background(), userCouponID, order.ID)
		}
	}()

	if l.svcCtx.UserRPC == nil {
		_ = l.svcCtx.Repo.UpdateStatus(ctx, orderNo, model.OrderStatusFailed)
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
		}
		return nil, errors.New("钱包服务不可用")
	}
	if err := l.svcCtx.UserRPC.Freeze(ctx, userID, payAmount, order.ID, orderNo); err != nil {
		_ = l.svcCtx.Repo.UpdateStatus(ctx, orderNo, model.OrderStatusFailed)
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
		}
		return nil, err
	}
	frozen := true
	defer func() {
		if frozen {
			_ = l.svcCtx.UserRPC.Unfreeze(context.Background(), userID, payAmount, order.ID, orderNo)
		}
	}()

	if err := cache.StockDeduct(ctx, l.svcCtx.Redis, redisItems); err != nil {
		_ = l.svcCtx.Repo.UpdateStatus(ctx, orderNo, model.OrderStatusFailed)
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
		}
		if errors.Is(err, cache.ErrStockInsufficient) || errors.Is(err, cache.ErrRedisUnavailable) {
			return nil, errors.New("库存不足")
		}
		return nil, fmt.Errorf("扣减库存失败: %w", err)
	}
	deducted := true
	defer func() {
		if deducted {
			_ = cache.StockRestore(context.Background(), l.svcCtx.Redis, redisItems)
			for _, h := range seckillHolds {
				_ = l.svcCtx.MerchantRPC.Restore(context.Background(), h.entryID, h.qty)
			}
		}
	}()

	if l.svcCtx.MQ == nil {
		_ = l.svcCtx.Repo.UpdateStatus(ctx, orderNo, model.OrderStatusFailed)
		return nil, errors.New("消息队列不可用")
	}
	if err := l.svcCtx.MQ.PublishOrderCreated(ctx, orderNo, stockItems); err != nil {
		_ = l.svcCtx.Repo.UpdateStatus(ctx, orderNo, model.OrderStatusFailed)
		return nil, fmt.Errorf("发布订单事件失败: %w", err)
	}

	frozen = false
	deducted = false
	couponLocked = false
	order.Items = orderItems
	return order, nil
}

func (l *OrderLogic) ListOrders(ctx context.Context, userID uint64, page, pageSize int, status string) ([]model.Order, int64, error) {
	orders, total, err := l.svcCtx.Repo.List(ctx, repository.OrderListFilter{
		UserID: userID, Page: page, PageSize: pageSize, Status: status,
	})
	if err != nil {
		return nil, 0, err
	}
	l.svcCtx.Repo.EnrichOrders(ctx, orders)
	return orders, total, nil
}

func (l *OrderLogic) UserOrderStatusCounts(ctx context.Context, userID uint64) (map[string]int64, error) {
	rows, err := l.svcCtx.Repo.CountByUserStatus(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := map[string]int64{
		model.OrderStatusPending:   0,
		model.OrderStatusConfirmed: 0,
		model.OrderStatusShipped:   0,
		model.OrderStatusCompleted: 0,
		model.OrderStatusReviewed:  0,
		model.OrderStatusCancelled: 0,
		model.OrderStatusFailed:    0,
		"after_sale":               0,
	}
	for _, row := range rows {
		out[row.Status] = row.Count
	}
	n, err := l.svcCtx.Repo.CountOpenAfterSalesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out["after_sale"] = n
	return out, nil
}

func (l *OrderLogic) ListUserAfterSales(ctx context.Context, userID uint64, page, pageSize int) ([]model.OrderAfterSale, int64, error) {
	return l.svcCtx.Repo.ListAfterSales(ctx, repository.AfterSaleListFilter{
		UserID: userID, Page: page, PageSize: pageSize,
	})
}

func (l *OrderLogic) GetOrder(ctx context.Context, userID, orderID uint64) (*model.Order, error) {
	order, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Repo.EnrichOrder(ctx, order)
	return order, nil
}

func (l *OrderLogic) notifyOrder(ctx context.Context, order *model.Order, title, content string) {
	if order == nil || l.svcCtx.UserRPC == nil || order.UserID == 0 {
		return
	}
	extra, _ := json.Marshal(map[string]interface{}{"order_no": order.OrderNo})
	_ = l.svcCtx.UserRPC.Notify(ctx, userrpc.NotifyReq{
		UserID: order.UserID, Title: title, Content: content,
		MsgType: "order", LinkType: "order", LinkID: order.ID, Extra: string(extra),
	})
}

func (l *OrderLogic) CancelOrder(ctx context.Context, userID, orderID uint64) error {
	order, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return errors.New("当前状态不可取消，请走售后")
	}
	if err := l.svcCtx.Repo.Cancel(ctx, orderID, userID); err != nil {
		return err
	}
	err = l.releaseStock(ctx, order)
	l.notifyOrder(ctx, order, "订单已取消", fmt.Sprintf("您的订单 %s 已取消", order.OrderNo))
	return err
}

func (l *OrderLogic) releaseStock(ctx context.Context, order *model.Order) error {
	items := make([]model.StockItem, 0, len(order.Items))
	redisItems := make([]cache.StockItem, 0, len(order.Items))
	for _, it := range order.Items {
		items = append(items, model.StockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
		if it.SkuID > 0 {
			redisItems = append(redisItems, cache.StockItem{SkuID: it.SkuID, Quantity: it.Quantity})
		}
		if it.SeckillEntryID > 0 && l.svcCtx.MerchantRPC != nil {
			_ = l.svcCtx.MerchantRPC.Restore(ctx, it.SeckillEntryID, it.Quantity)
		}
	}
	if l.svcCtx.Redis != nil && len(redisItems) > 0 {
		_ = cache.StockRestore(ctx, l.svcCtx.Redis, redisItems)
	}
	pay := orderPayAmount(order)
	if l.svcCtx.UserRPC != nil && pay > 0 {
		_ = l.svcCtx.UserRPC.Unfreeze(ctx, order.UserID, pay, order.ID, order.OrderNo)
	}
	if order.UserCouponID > 0 && l.svcCtx.MerchantRPC != nil {
		_ = l.svcCtx.MerchantRPC.UnlockCoupon(ctx, order.UserCouponID, order.ID)
	}
	if l.svcCtx.MQ != nil {
		_ = l.svcCtx.MQ.PublishOrderCancelled(ctx, order.OrderNo, items)
	}
	return nil
}

func (l *OrderLogic) ListFiltered(ctx context.Context, f repository.OrderListFilter) ([]model.Order, int64, error) {
	orders, total, err := l.svcCtx.Repo.List(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	l.svcCtx.Repo.EnrichOrders(ctx, orders)
	return orders, total, nil
}

func (l *OrderLogic) ListByShop(ctx context.Context, shopID uint64, page, pageSize int, status, orderNo string) ([]model.Order, int64, error) {
	return l.ListFiltered(ctx, repository.OrderListFilter{
		ShopID: shopID, Page: page, PageSize: pageSize, Status: status, OrderNo: orderNo,
	})
}

func (l *OrderLogic) ListAll(ctx context.Context, shopID uint64, page, pageSize int, status, orderNo string) ([]model.Order, int64, error) {
	return l.ListFiltered(ctx, repository.OrderListFilter{
		ShopID: shopID, Page: page, PageSize: pageSize, Status: status, OrderNo: orderNo,
	})
}

func (l *OrderLogic) GetOrderByShop(ctx context.Context, shopID, orderID uint64) (*model.Order, error) {
	order, err := l.svcCtx.Repo.FindByIDAndShop(ctx, orderID, shopID)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Repo.EnrichOrder(ctx, order)
	return order, nil
}

func (l *OrderLogic) GetOrderAdmin(ctx context.Context, orderID uint64) (*model.Order, error) {
	order, err := l.svcCtx.Repo.FindByIDAdmin(ctx, orderID)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Repo.EnrichOrder(ctx, order)
	return order, nil
}

func (l *OrderLogic) Ship(ctx context.Context, id, shopID uint64, company, shipNo string) error {
	if company == "" || shipNo == "" {
		return errors.New("物流公司与单号必填")
	}
	var order *model.Order
	var err error
	if shopID > 0 {
		order, err = l.svcCtx.Repo.FindByIDAndShop(ctx, id, shopID)
	} else {
		order, err = l.svcCtx.Repo.FindByIDAdmin(ctx, id)
	}
	if err != nil {
		return errors.New("订单不存在或状态不是待发货(confirmed)")
	}
	if err := l.svcCtx.Repo.Ship(ctx, id, shopID, company, shipNo); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return errors.New("订单不存在或状态不是待发货(confirmed)")
		}
		return err
	}
	l.notifyOrder(ctx, order, "商家已发货", fmt.Sprintf("订单 %s 已发货，物流：%s %s", order.OrderNo, company, shipNo))
	return nil
}

func (l *OrderLogic) Complete(ctx context.Context, id, shopID uint64) error {
	var order *model.Order
	var err error
	if shopID > 0 {
		order, err = l.svcCtx.Repo.FindByIDAndShop(ctx, id, shopID)
	} else {
		order, err = l.svcCtx.Repo.FindByIDAdmin(ctx, id)
	}
	if err != nil {
		return errors.New("订单不存在或状态不是已发货")
	}
	if err := l.svcCtx.Repo.Complete(ctx, id, shopID); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return errors.New("订单不存在或状态不是已发货")
		}
		return err
	}
	l.redeemCoupon(ctx, order)
	l.notifyOrder(ctx, order, "订单已完成", fmt.Sprintf("订单 %s 已完成，欢迎再次光临", order.OrderNo))
	return nil
}

func (l *OrderLogic) ConfirmReceive(ctx context.Context, userID, orderID uint64) error {
	order, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)
	if err != nil {
		return errors.New("订单不存在或状态不是已发货")
	}
	if err := l.svcCtx.Repo.ConfirmReceive(ctx, orderID, userID); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return errors.New("订单不存在或状态不是已发货")
		}
		return err
	}
	l.redeemCoupon(ctx, order)
	l.notifyOrder(ctx, order, "订单已完成", fmt.Sprintf("您已确认收货，订单 %s 已完成", order.OrderNo))
	return nil
}

func (l *OrderLogic) redeemCoupon(ctx context.Context, order *model.Order) {
	if order == nil || order.UserCouponID == 0 || l.svcCtx.MerchantRPC == nil {
		return
	}
	_ = l.svcCtx.MerchantRPC.RedeemCoupon(ctx, order.UserCouponID, order.ID, order.DiscountAmount)
}

func (l *OrderLogic) UpdateRemark(ctx context.Context, id, shopID uint64, remark string) error {
	if err := l.svcCtx.Repo.UpdateRemark(ctx, id, shopID, remark); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	return nil
}

func (l *OrderLogic) CreateAfterSale(ctx context.Context, userID, orderID uint64, req types.CreateAfterSaleReq) (*model.OrderAfterSale, error) {
	order, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	gate := l.evaluateAfterSale(ctx, order)
	if !gate.Eligible {
		return nil, errors.New(gate.Reason)
	}
	n, _ := l.svcCtx.Repo.CountOpenAfterSales(ctx, orderID)
	if n > 0 {
		return nil, errors.New("已有进行中的售后单")
	}
	typ := req.Type
	if typ == "" {
		typ = model.AfterSaleTypeRefund
	}
	if typ != model.AfterSaleTypeRefund && typ != model.AfterSaleTypeReturnRefund {
		return nil, errors.New("售后类型无效")
	}
	amount := req.Amount
	pay := orderPayAmount(order)
	if amount <= 0 {
		amount = pay
	}
	if amount > pay {
		return nil, errors.New("退款金额不能超过订单金额")
	}
	as := &model.OrderAfterSale{
		OrderID: order.ID,
		OrderNo: order.OrderNo,
		UserID:  order.UserID,
		ShopID:  order.ShopID,
		Type:    typ,
		Reason:  req.Reason,
		Amount:  amount,
		Status:  model.AfterSalePending,
	}
	if err := l.svcCtx.Repo.CreateAfterSale(ctx, as); err != nil {
		return nil, err
	}
	return as, nil
}

func (l *OrderLogic) AfterSaleEligible(ctx context.Context, userID, orderID uint64) (*types.AfterSaleEligibleResp, error) {
	order, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	gate := l.evaluateAfterSale(ctx, order)
	if gate.Eligible {
		n, _ := l.svcCtx.Repo.CountOpenAfterSales(ctx, orderID)
		if n > 0 {
			gate.Eligible = false
			gate.Reason = "已有进行中的售后单"
		}
	}
	return &types.AfterSaleEligibleResp{
		Eligible: gate.Eligible,
		Reason:   gate.Reason,
		Days:     gate.Days,
		Deadline: gate.Deadline,
	}, nil
}

const (
	configKeyAfterSaleDays = "order_after_sale_days"
	defaultAfterSaleDays   = 7
)

type afterSaleGate struct {
	Eligible bool
	Reason   string
	Days     int
	Deadline string
}

func (l *OrderLogic) afterSaleWindowDays(ctx context.Context) int {
	if l.svcCtx.UserRPC == nil {
		return defaultAfterSaleDays
	}
	val, err := l.svcCtx.UserRPC.GetConfig(ctx, configKeyAfterSaleDays)
	if err != nil || strings.TrimSpace(val) == "" {
		return defaultAfterSaleDays
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil || n < 0 {
		return defaultAfterSaleDays
	}
	return n
}

func (l *OrderLogic) evaluateAfterSale(ctx context.Context, order *model.Order) afterSaleGate {
	days := l.afterSaleWindowDays(ctx)
	out := afterSaleGate{Days: days}
	switch order.Status {
	case model.OrderStatusConfirmed, model.OrderStatusShipped, model.OrderStatusCompleted, model.OrderStatusReviewed:
	default:
		out.Reason = "当前订单状态不可申请售后"
		return out
	}
	if order.CompletedAt != nil {
		completed := time.Time(*order.CompletedAt)
		deadline := completed.Add(time.Duration(days) * 24 * time.Hour)
		out.Deadline = deadline.Format("2006-01-02 15:04:05")
		if time.Now().After(deadline) {
			out.Reason = "已超过售后申请期限"
			return out
		}
	}
	out.Eligible = true
	return out
}

func (l *OrderLogic) ListAfterSales(ctx context.Context, f repository.AfterSaleListFilter) ([]model.OrderAfterSale, int64, error) {
	return l.svcCtx.Repo.ListAfterSales(ctx, f)
}

func (l *OrderLogic) ListAfterSalesByOrder(ctx context.Context, orderID uint64) ([]model.OrderAfterSale, error) {
	return l.svcCtx.Repo.ListAfterSalesByOrder(ctx, orderID)
}

func (l *OrderLogic) HandleAfterSale(ctx context.Context, id, shopID, handledBy uint64, action, adminRemark string) error {
	as, err := l.svcCtx.Repo.FindAfterSale(ctx, id)
	if err != nil {
		return errors.New("售后单不存在")
	}
	if shopID > 0 && as.ShopID != shopID {
		return errors.New("无权处理该售后单")
	}
	var status string
	switch action {
	case "approve":
		if as.Status != model.AfterSalePending {
			return errors.New("仅待处理可同意")
		}
		status = model.AfterSaleApproved
	case "reject":
		if as.Status != model.AfterSalePending {
			return errors.New("仅待处理可拒绝")
		}
		status = model.AfterSaleRejected
	case "refunded":
		if as.Status != model.AfterSalePending && as.Status != model.AfterSaleApproved {
			return errors.New("当前状态不可退款完成")
		}
		status = model.AfterSaleRefunded
	case "closed":
		status = model.AfterSaleClosed
	default:
		return errors.New("无效操作")
	}
	if err := l.svcCtx.Repo.HandleAfterSale(ctx, id, shopID, handledBy, status, adminRemark); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return errors.New("售后单不存在或无权处理")
		}
		return err
	}
	if status == model.AfterSaleRefunded {
		order, err := l.svcCtx.Repo.FindByIDAdmin(ctx, as.OrderID)
		if err != nil {
			return nil
		}
		pay := orderPayAmount(order)
		// 全额退：返还优惠券（确认收货后券已 used，走 return）
		if as.Amount+0.001 >= pay && order.UserCouponID > 0 && l.svcCtx.MerchantRPC != nil {
			_ = l.svcCtx.MerchantRPC.ReturnCoupon(ctx, order.UserCouponID, order.ID)
		}
		// 已扣库存的状态：还库存；并标记订单取消（若尚未取消）
		if order.Status == model.OrderStatusConfirmed || order.Status == model.OrderStatusShipped ||
			order.Status == model.OrderStatusCompleted || order.Status == model.OrderStatusReviewed {
			// releaseStock 内 unlock 对已 used 券也会走 return 语义前的 unlock；全额已单独 Return
			if as.Amount+0.001 < pay {
				// 部分退不 unlock（券保持 used）；临时清 user_coupon 避免 release 解锁
				uc := order.UserCouponID
				order.UserCouponID = 0
				_ = l.releaseStock(ctx, order)
				order.UserCouponID = uc
			} else {
				order.UserCouponID = 0 // 已 Return，避免再 Unlock
				_ = l.releaseStock(ctx, order)
			}
			_ = l.svcCtx.Repo.UpdateStatus(ctx, order.OrderNo, model.OrderStatusCancelled)
		}
		l.notifyOrder(ctx, order, "退款已完成", fmt.Sprintf("订单 %s 退款已完成，金额 ¥%.2f", order.OrderNo, as.Amount))
	}
	return nil
}
