package repository

import (
	"mymall/services/order-service/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
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
	var total int64
	if err := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []model.Order
	offset := (page - 1) * pageSize
	err := r.db.Preload("Items").Where("user_id = ?", userID).
		Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
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
	var total int64
	if err := r.db.Model(&model.Order{}).Where("shop_id = ?", shopID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []model.Order
	offset := (page - 1) * pageSize
	err := r.db.Preload("Items").Where("shop_id = ?", shopID).
		Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) ListAll(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	q := r.db.Model(&model.Order{})
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query := r.db.Preload("Items")
	if shopID > 0 {
		query = query.Where("shop_id = ?", shopID)
	}
	var orders []model.Order
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
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
