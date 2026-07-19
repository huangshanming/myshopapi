package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mymall/pkg/cache"
	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderLogic struct {
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{svcCtx: svcCtx}
}

func (l *OrderLogic) CreateOrder(ctx context.Context, userID uint64, addressID uint64, items []types.CreateOrderItem) (*model.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}
	if addressID == 0 {
		return nil, errors.New("请选择收货地址")
	}
	if l.svcCtx.Redis == nil {
		return nil, errors.New("库存服务不可用，请稍后重试")
	}
	if l.svcCtx.UserHTTP == nil {
		return nil, errors.New("用户服务不可用")
	}
	addr, err := l.svcCtx.UserHTTP.GetAddress(ctx, userID, addressID)
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
			if l.svcCtx.MerchantHTTP == nil {
				return nil, errors.New("秒杀服务不可用")
			}
			cr, err := l.svcCtx.MerchantHTTP.Consume(ctx, it.SeckillEntryID, it.ProductID, qty)
			if err != nil {
				for _, h := range seckillHolds {
					_ = l.svcCtx.MerchantHTTP.Restore(context.Background(), h.entryID, h.qty)
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

	orderNo := fmt.Sprintf("ORD%s", uuid.NewString()[:8]+time.Now().Format("150405"))
	order := &model.Order{
		OrderNo:         orderNo,
		UserID:          userID,
		ShopID:          shopID,
		TotalAmount:     total,
		ReceiverName:    addr.ReceiverName,
		ReceiverPhone:   addr.ReceiverPhone,
		ReceiverAddress: addr.FullAddress(),
		Status:          model.OrderStatusPending,
	}
	if err := l.svcCtx.Repo.Create(order, orderItems); err != nil {
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantHTTP.Restore(context.Background(), h.entryID, h.qty)
		}
		return nil, err
	}

	if l.svcCtx.UserHTTP == nil {
		_ = l.svcCtx.Repo.UpdateStatus(orderNo, model.OrderStatusFailed)
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantHTTP.Restore(context.Background(), h.entryID, h.qty)
		}
		return nil, errors.New("钱包服务不可用")
	}
	if err := l.svcCtx.UserHTTP.Freeze(ctx, userID, total, order.ID, orderNo); err != nil {
		_ = l.svcCtx.Repo.UpdateStatus(orderNo, model.OrderStatusFailed)
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantHTTP.Restore(context.Background(), h.entryID, h.qty)
		}
		return nil, err
	}
	frozen := true
	defer func() {
		if frozen {
			_ = l.svcCtx.UserHTTP.Unfreeze(context.Background(), userID, total, order.ID, orderNo)
		}
	}()

	if err := cache.StockDeduct(ctx, l.svcCtx.Redis, redisItems); err != nil {
		_ = l.svcCtx.Repo.UpdateStatus(orderNo, model.OrderStatusFailed)
		for _, h := range seckillHolds {
			_ = l.svcCtx.MerchantHTTP.Restore(context.Background(), h.entryID, h.qty)
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
				_ = l.svcCtx.MerchantHTTP.Restore(context.Background(), h.entryID, h.qty)
			}
		}
	}()

	if l.svcCtx.MQ == nil {
		_ = l.svcCtx.Repo.UpdateStatus(orderNo, model.OrderStatusFailed)
		return nil, errors.New("消息队列不可用")
	}
	if err := l.svcCtx.MQ.PublishOrderCreated(ctx, orderNo, stockItems); err != nil {
		_ = l.svcCtx.Repo.UpdateStatus(orderNo, model.OrderStatusFailed)
		return nil, fmt.Errorf("发布订单事件失败: %w", err)
	}

	frozen = false
	deducted = false
	order.Items = orderItems
	return order, nil
}

func (l *OrderLogic) ListOrders(userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := l.svcCtx.Repo.List(repository.OrderListFilter{UserID: userID, Page: page, PageSize: pageSize})
	if err != nil {
		return nil, 0, err
	}
	l.svcCtx.Repo.EnrichOrders(orders)
	return orders, total, nil
}

func (l *OrderLogic) GetOrder(userID, orderID uint64) (*model.Order, error) {
	order, err := l.svcCtx.Repo.FindByID(orderID, userID)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Repo.EnrichOrder(order)
	return order, nil
}

func (l *OrderLogic) CancelOrder(ctx context.Context, userID, orderID uint64) error {
	order, err := l.svcCtx.Repo.FindByID(orderID, userID)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return errors.New("当前状态不可取消，请走售后")
	}
	if err := l.svcCtx.Repo.Cancel(orderID, userID); err != nil {
		return err
	}
	return l.releaseStock(ctx, order)
}

func (l *OrderLogic) releaseStock(ctx context.Context, order *model.Order) error {
	items := make([]model.StockItem, 0, len(order.Items))
	redisItems := make([]cache.StockItem, 0, len(order.Items))
	for _, it := range order.Items {
		items = append(items, model.StockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
		if it.SkuID > 0 {
			redisItems = append(redisItems, cache.StockItem{SkuID: it.SkuID, Quantity: it.Quantity})
		}
		if it.SeckillEntryID > 0 && l.svcCtx.MerchantHTTP != nil {
			_ = l.svcCtx.MerchantHTTP.Restore(ctx, it.SeckillEntryID, it.Quantity)
		}
	}
	if l.svcCtx.Redis != nil && len(redisItems) > 0 {
		_ = cache.StockRestore(ctx, l.svcCtx.Redis, redisItems)
	}
	if l.svcCtx.UserHTTP != nil && order.TotalAmount > 0 {
		_ = l.svcCtx.UserHTTP.Unfreeze(ctx, order.UserID, order.TotalAmount, order.ID, order.OrderNo)
	}
	if l.svcCtx.MQ != nil {
		_ = l.svcCtx.MQ.PublishOrderCancelled(ctx, order.OrderNo, items)
	}
	return nil
}

func (l *OrderLogic) ListFiltered(f repository.OrderListFilter) ([]model.Order, int64, error) {
	orders, total, err := l.svcCtx.Repo.List(f)
	if err != nil {
		return nil, 0, err
	}
	l.svcCtx.Repo.EnrichOrders(orders)
	return orders, total, nil
}

func (l *OrderLogic) ListByShop(shopID uint64, page, pageSize int, status, orderNo string) ([]model.Order, int64, error) {
	return l.ListFiltered(repository.OrderListFilter{
		ShopID: shopID, Page: page, PageSize: pageSize, Status: status, OrderNo: orderNo,
	})
}

func (l *OrderLogic) ListAll(shopID uint64, page, pageSize int, status, orderNo string) ([]model.Order, int64, error) {
	return l.ListFiltered(repository.OrderListFilter{
		ShopID: shopID, Page: page, PageSize: pageSize, Status: status, OrderNo: orderNo,
	})
}

func (l *OrderLogic) GetOrderByShop(shopID, orderID uint64) (*model.Order, error) {
	order, err := l.svcCtx.Repo.FindByIDAndShop(orderID, shopID)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Repo.EnrichOrder(order)
	return order, nil
}

func (l *OrderLogic) GetOrderAdmin(orderID uint64) (*model.Order, error) {
	order, err := l.svcCtx.Repo.FindByIDAdmin(orderID)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Repo.EnrichOrder(order)
	return order, nil
}

func (l *OrderLogic) Ship(id, shopID uint64, company, shipNo string) error {
	if company == "" || shipNo == "" {
		return errors.New("物流公司与单号必填")
	}
	if err := l.svcCtx.Repo.Ship(id, shopID, company, shipNo); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在或状态不是待发货(confirmed)")
		}
		return err
	}
	return nil
}

func (l *OrderLogic) Complete(id, shopID uint64) error {
	if err := l.svcCtx.Repo.Complete(id, shopID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在或状态不是已发货")
		}
		return err
	}
	return nil
}

func (l *OrderLogic) UpdateRemark(id, shopID uint64, remark string) error {
	if err := l.svcCtx.Repo.UpdateRemark(id, shopID, remark); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	return nil
}

func (l *OrderLogic) CreateAfterSale(userID, orderID uint64, req types.CreateAfterSaleReq) (*model.OrderAfterSale, error) {
	order, err := l.svcCtx.Repo.FindByID(orderID, userID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	switch order.Status {
	case model.OrderStatusConfirmed, model.OrderStatusShipped, model.OrderStatusCompleted:
	default:
		return nil, errors.New("当前订单状态不可申请售后")
	}
	n, _ := l.svcCtx.Repo.CountOpenAfterSales(orderID)
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
	if amount <= 0 {
		amount = order.TotalAmount
	}
	if amount > order.TotalAmount {
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
	if err := l.svcCtx.Repo.CreateAfterSale(as); err != nil {
		return nil, err
	}
	return as, nil
}

func (l *OrderLogic) ListAfterSales(f repository.AfterSaleListFilter) ([]model.OrderAfterSale, int64, error) {
	return l.svcCtx.Repo.ListAfterSales(f)
}

func (l *OrderLogic) ListAfterSalesByOrder(orderID uint64) ([]model.OrderAfterSale, error) {
	return l.svcCtx.Repo.ListAfterSalesByOrder(orderID)
}

func (l *OrderLogic) HandleAfterSale(ctx context.Context, id, shopID, handledBy uint64, action, adminRemark string) error {
	as, err := l.svcCtx.Repo.FindAfterSale(id)
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
	if err := l.svcCtx.Repo.HandleAfterSale(id, shopID, handledBy, status, adminRemark); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("售后单不存在或无权处理")
		}
		return err
	}
	if status == model.AfterSaleRefunded {
		order, err := l.svcCtx.Repo.FindByIDAdmin(as.OrderID)
		if err != nil {
			return nil
		}
		// 已扣库存的状态：还库存；并标记订单取消（若尚未取消）
		if order.Status == model.OrderStatusConfirmed || order.Status == model.OrderStatusShipped || order.Status == model.OrderStatusCompleted {
			_ = l.releaseStock(ctx, order)
			_ = l.svcCtx.Repo.UpdateStatus(order.OrderNo, model.OrderStatusCancelled)
		}
	}
	return nil
}
