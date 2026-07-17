package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/google/uuid"
)

type OrderLogic struct {
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{svcCtx: svcCtx}
}

func (l *OrderLogic) CreateOrder(ctx context.Context, userID uint64, items []types.CreateOrderItem) (*model.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("订单商品不能为空")
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
		stock      int32
		defSku     uint64
	}, len(resp.GetProducts()))
	for _, p := range resp.GetProducts() {
		productByID[p.GetId()] = &struct {
			id, shopID uint64
			name       string
			price      float64
			stock      int32
			defSku     uint64
		}{p.GetId(), p.GetShopId(), p.GetName(), p.GetSalePrice(), p.GetStock(), p.GetDefaultSkuId()}
	}

	var total float64
	var shopID uint64
	orderItems := make([]model.OrderItem, 0, len(items))
	stockItems := make([]model.StockItem, 0, len(items))
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
		qty := it.Quantity
		if p.stock < int32(qty) {
			return nil, fmt.Errorf("商品 %s 库存不足", p.name)
		}
		snap := it.SkuSnapshot
		if snap == "" && skuID > 0 {
			snap = fmt.Sprintf(`{"sku_id":%d}`, skuID)
		}
		lineAmount := p.price * float64(qty)
		total += lineAmount
		orderItems = append(orderItems, model.OrderItem{
			ProductID:   p.id,
			SkuID:       skuID,
			ProductName: p.name,
			SkuSnapshot: snap,
			Price:       p.price,
			Quantity:    qty,
		})
		stockItems = append(stockItems, model.StockItem{ProductID: p.id, SkuID: skuID, Quantity: qty})
	}

	orderNo := fmt.Sprintf("ORD%s", uuid.NewString()[:8]+time.Now().Format("150405"))
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		ShopID:      shopID,
		TotalAmount: total,
		Status:      model.OrderStatusPending,
	}
	if err := l.svcCtx.Repo.Create(order, orderItems); err != nil {
		return nil, err
	}

	if l.svcCtx.MQ != nil {
		if err := l.svcCtx.MQ.PublishOrderCreated(ctx, orderNo, stockItems); err != nil {
			_ = l.svcCtx.Repo.UpdateStatus(orderNo, model.OrderStatusFailed)
			return nil, fmt.Errorf("发布订单事件失败: %w", err)
		}
	}

	order.Items = orderItems
	return order, nil
}

func (l *OrderLogic) ListOrders(userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListByUser(userID, page, pageSize)
}

func (l *OrderLogic) GetOrder(userID, orderID uint64) (*model.Order, error) {
	return l.svcCtx.Repo.FindByID(orderID, userID)
}

func (l *OrderLogic) CancelOrder(ctx context.Context, userID, orderID uint64) error {
	order, err := l.svcCtx.Repo.FindByID(orderID, userID)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return errors.New("当前状态不可取消")
	}
	if err := l.svcCtx.Repo.Cancel(orderID, userID); err != nil {
		return err
	}
	if l.svcCtx.MQ != nil {
		items := make([]model.StockItem, 0, len(order.Items))
		for _, it := range order.Items {
			items = append(items, model.StockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
		}
		_ = l.svcCtx.MQ.PublishOrderCancelled(ctx, order.OrderNo, items)
	}
	return nil
}

func (l *OrderLogic) ListByShop(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListByShop(shopID, page, pageSize)
}

func (l *OrderLogic) ListAll(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListAll(shopID, page, pageSize)
}

func (l *OrderLogic) GetOrderByShop(shopID, orderID uint64) (*model.Order, error) {
	return l.svcCtx.Repo.FindByIDAndShop(orderID, shopID)
}

func (l *OrderLogic) GetOrderAdmin(orderID uint64) (*model.Order, error) {
	return l.svcCtx.Repo.FindByIDAdmin(orderID)
}
