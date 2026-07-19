package repository

import (
	"time"

	"mymall/common"
	"mymall/services/order-service/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type OrderListFilter struct {
	ShopID  uint64
	UserID  uint64
	Status  string
	OrderNo string
	Page    int
	PageSize int
}

func (r *OrderRepository) Create(order *model.Order, items []model.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
		}
		return tx.Create(&items).Error
	})
}

func (r *OrderRepository) FindByID(id, userID uint64) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("id = ? AND user_id = ?", id, userID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListByUser(userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return r.List(OrderListFilter{UserID: userID, Page: page, PageSize: pageSize})
}

func (r *OrderRepository) UpdateStatus(orderNo, status string) error {
	return r.db.Model(&model.Order{}).Where("order_no = ?", orderNo).Update("status", status).Error
}

func (r *OrderRepository) FindByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Cancel(orderID, userID uint64) error {
	return r.db.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND status IN ?", orderID, userID, []string{model.OrderStatusPending, model.OrderStatusConfirmed}).
		Update("status", model.OrderStatusCancelled).Error
}

func (r *OrderRepository) ListByShop(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return r.List(OrderListFilter{ShopID: shopID, Page: page, PageSize: pageSize})
}

func (r *OrderRepository) ListAll(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return r.List(OrderListFilter{ShopID: shopID, Page: page, PageSize: pageSize})
}

func (r *OrderRepository) List(f OrderListFilter) ([]model.Order, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	q := r.db.Model(&model.Order{})
	if f.ShopID > 0 {
		q = q.Where("shop_id = ?", f.ShopID)
	}
	if f.UserID > 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.OrderNo != "" {
		q = q.Where("order_no LIKE ?", "%"+f.OrderNo+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []model.Order
	err := q.Preload("Items").Order("id DESC").
		Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) FindByIDAndShop(id, shopID uint64) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("id = ? AND shop_id = ?", id, shopID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByIDAdmin(id uint64) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Ship(id uint64, shopID uint64, company, shipNo string) error {
	now := common.LocalTime(time.Now())
	q := r.db.Model(&model.Order{}).Where("id = ? AND status = ?", id, model.OrderStatusConfirmed)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Updates(map[string]interface{}{
		"status":       model.OrderStatusShipped,
		"ship_company": company,
		"ship_no":      shipNo,
		"shipped_at":   &now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) Complete(id uint64, shopID uint64) error {
	now := common.LocalTime(time.Now())
	q := r.db.Model(&model.Order{}).Where("id = ? AND status = ?", id, model.OrderStatusShipped)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Updates(map[string]interface{}{
		"status":       model.OrderStatusCompleted,
		"completed_at": &now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) ConfirmReceive(id, userID uint64) error {
	now := common.LocalTime(time.Now())
	res := r.db.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND status = ?", id, userID, model.OrderStatusShipped).
		Updates(map[string]interface{}{
			"status":       model.OrderStatusCompleted,
			"completed_at": &now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) MarkReviewed(id, userID uint64) error {
	now := common.LocalTime(time.Now())
	res := r.db.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND status = ?", id, userID, model.OrderStatusCompleted).
		Updates(map[string]interface{}{
			"status":      model.OrderStatusReviewed,
			"reviewed_at": &now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) UpdateRemark(id uint64, shopID uint64, remark string) error {
	q := r.db.Model(&model.Order{}).Where("id = ?", id)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Update("remark", remark)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) EnrichOrders(orders []model.Order) {
	if len(orders) == 0 {
		return
	}
	userIDs := make([]uint64, 0, len(orders))
	shopIDs := make([]uint64, 0, len(orders))
	seenU, seenS := map[uint64]struct{}{}, map[uint64]struct{}{}
	for _, o := range orders {
		if _, ok := seenU[o.UserID]; !ok && o.UserID > 0 {
			seenU[o.UserID] = struct{}{}
			userIDs = append(userIDs, o.UserID)
		}
		if _, ok := seenS[o.ShopID]; !ok && o.ShopID > 0 {
			seenS[o.ShopID] = struct{}{}
			shopIDs = append(shopIDs, o.ShopID)
		}
	}
	userNames := r.loadUserNames(userIDs)
	shopNames := r.loadShopNames(shopIDs)
	for i := range orders {
		orders[i].UserName = userNames[orders[i].UserID]
		orders[i].ShopName = shopNames[orders[i].ShopID]
	}
}

func (r *OrderRepository) EnrichOrder(order *model.Order) {
	if order == nil {
		return
	}
	r.EnrichOrders([]model.Order{*order})
	// EnrichOrders works on copy — fix by direct load
	names := r.loadUserNames([]uint64{order.UserID})
	shops := r.loadShopNames([]uint64{order.ShopID})
	order.UserName = names[order.UserID]
	order.ShopName = shops[order.ShopID]
}

func (r *OrderRepository) loadUserNames(ids []uint64) map[uint64]string {
	out := map[uint64]string{}
	if len(ids) == 0 {
		return out
	}
	type row struct {
		ID       uint64 `gorm:"column:id"`
		Nickname string `gorm:"column:nickname"`
		Mobile   string `gorm:"column:mobile"`
	}
	var rows []row
	if err := r.db.Table("users").Select("id, nickname, mobile").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return out
	}
	for _, u := range rows {
		if u.Nickname != "" {
			out[u.ID] = u.Nickname
		} else if u.Mobile != "" {
			out[u.ID] = u.Mobile
		}
	}
	return out
}

func (r *OrderRepository) loadShopNames(ids []uint64) map[uint64]string {
	out := map[uint64]string{}
	if len(ids) == 0 {
		return out
	}
	// shops 表字段为 name（非 shop_name）
	type row struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var rows []row
	if err := r.db.Table("shops").Select("id, name").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return out
	}
	for _, s := range rows {
		out[s.ID] = s.Name
	}
	return out
}

// After-sales

type AfterSaleListFilter struct {
	ShopID   uint64
	UserID   uint64
	Status   string
	OrderNo  string
	Page     int
	PageSize int
}

func (r *OrderRepository) CreateAfterSale(as *model.OrderAfterSale) error {
	return r.db.Create(as).Error
}

func (r *OrderRepository) FindAfterSale(id uint64) (*model.OrderAfterSale, error) {
	var as model.OrderAfterSale
	if err := r.db.Where("id = ?", id).First(&as).Error; err != nil {
		return nil, err
	}
	return &as, nil
}

func (r *OrderRepository) ListAfterSales(f AfterSaleListFilter) ([]model.OrderAfterSale, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	q := r.db.Model(&model.OrderAfterSale{})
	if f.ShopID > 0 {
		q = q.Where("shop_id = ?", f.ShopID)
	}
	if f.UserID > 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.OrderNo != "" {
		q = q.Where("order_no LIKE ?", "%"+f.OrderNo+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.OrderAfterSale
	err := q.Order("id DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	r.enrichAfterSales(list)
	return list, total, nil
}

func (r *OrderRepository) ListAfterSalesByOrder(orderID uint64) ([]model.OrderAfterSale, error) {
	var list []model.OrderAfterSale
	err := r.db.Where("order_id = ?", orderID).Order("id DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	r.enrichAfterSales(list)
	return list, nil
}

func (r *OrderRepository) HandleAfterSale(id, shopID, handledBy uint64, status, adminRemark string) error {
	q := r.db.Model(&model.OrderAfterSale{}).Where("id = ?", id)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Updates(map[string]interface{}{
		"status":       status,
		"admin_remark": adminRemark,
		"handled_by":   handledBy,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) CountOpenAfterSales(orderID uint64) (int64, error) {
	var n int64
	err := r.db.Model(&model.OrderAfterSale{}).
		Where("order_id = ? AND status IN ?", orderID, []string{model.AfterSalePending, model.AfterSaleApproved}).
		Count(&n).Error
	return n, err
}

func (r *OrderRepository) enrichAfterSales(list []model.OrderAfterSale) {
	if len(list) == 0 {
		return
	}
	userIDs, shopIDs := []uint64{}, []uint64{}
	seenU, seenS := map[uint64]struct{}{}, map[uint64]struct{}{}
	for _, a := range list {
		if _, ok := seenU[a.UserID]; !ok {
			seenU[a.UserID] = struct{}{}
			userIDs = append(userIDs, a.UserID)
		}
		if _, ok := seenS[a.ShopID]; !ok && a.ShopID > 0 {
			seenS[a.ShopID] = struct{}{}
			shopIDs = append(shopIDs, a.ShopID)
		}
	}
	un := r.loadUserNames(userIDs)
	sn := r.loadShopNames(shopIDs)
	for i := range list {
		list[i].UserName = un[list[i].UserID]
		list[i].ShopName = sn[list[i].ShopID]
	}
}
