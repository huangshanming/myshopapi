package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PointsOrderRepository struct {
	db *gorm.DB
}

func NewPointsOrderRepository(db *gorm.DB) *PointsOrderRepository {
	return &PointsOrderRepository{db: db}
}

type PointsOrderListFilter struct {
	Status  string
	OrderNo string
	UserID  uint64
	Keyword string
}

func (r *PointsOrderRepository) List(page, pageSize int, f PointsOrderListFilter) ([]model.PointsExchangeOrder, int64, error) {
	q := r.db.Model(&model.PointsExchangeOrder{})
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.OrderNo != "" {
		q = q.Where("order_no = ?", strings.TrimSpace(f.OrderNo))
	}
	if f.UserID > 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		q = q.Where("product_name LIKE ? OR order_no LIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.PointsExchangeOrder
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *PointsOrderRepository) GetByID(id uint64) (*model.PointsExchangeOrder, error) {
	var o model.PointsExchangeOrder
	if err := r.db.First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func genOrderNo() string {
	return fmt.Sprintf("PE%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
}

// CreateExchangeLocal 本地事务：校验商品、扣库存、创建兑换单（积分扣减由调用方走 user-service）
func (r *PointsOrderRepository) CreateExchangeLocal(userID, productID uint64, quantity int, receiverName, receiverPhone, receiverAddress string) (*model.PointsExchangeOrder, error) {
	if quantity < 1 {
		quantity = 1
	}
	var out *model.PointsExchangeOrder
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var p model.PointsProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, productID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("商品不存在")
			}
			return err
		}
		if p.Status != model.PointsProductStatusOn {
			return errors.New("商品已下架")
		}
		if p.Stock < quantity {
			return errors.New("库存不足")
		}
		cost := p.PointsPrice * quantity
		if cost <= 0 {
			return errors.New("商品积分价无效")
		}
		if p.PerUserLimit > 0 {
			var n int64
			if err := tx.Model(&model.PointsExchangeOrder{}).
				Where("user_id = ? AND product_id = ? AND status <> ?", userID, productID, model.PointsOrderCancelled).
				Count(&n).Error; err != nil {
				return err
			}
			if int(n)+quantity > p.PerUserLimit {
				return errors.New("已达兑换上限")
			}
		}
		if err := tx.Model(&p).Update("stock", p.Stock-quantity).Error; err != nil {
			return err
		}
		o := &model.PointsExchangeOrder{
			OrderNo:         genOrderNo(),
			UserID:          userID,
			ProductID:       p.ID,
			ProductName:     p.Name,
			ProductCover:    p.CoverURL,
			Quantity:        quantity,
			PointsCost:      cost,
			Status:          model.PointsOrderPending,
			ReceiverName:    strings.TrimSpace(receiverName),
			ReceiverPhone:   strings.TrimSpace(receiverPhone),
			ReceiverAddress: strings.TrimSpace(receiverAddress),
		}
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		out = o
		return nil
	})
	return out, err
}

// AbortExchange 扣积分失败时回滚：删单并退库存
func (r *PointsOrderRepository) AbortExchange(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var o model.PointsExchangeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&o, id).Error; err != nil {
			return err
		}
		if o.Status != model.PointsOrderPending {
			return nil
		}
		if err := tx.Model(&model.PointsProduct{}).Where("id = ?", o.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", o.Quantity)).Error; err != nil {
			return err
		}
		return tx.Delete(&o).Error
	})
}

func (r *PointsOrderRepository) Ship(id uint64, company, shipNo string) (*model.PointsExchangeOrder, error) {
	o, err := r.GetByID(id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if o.Status != model.PointsOrderPending {
		return nil, errors.New("当前状态不可发货")
	}
	now := common.LocalTime(time.Now())
	if err := r.db.Model(o).Updates(map[string]interface{}{
		"status": model.PointsOrderShipped, "ship_company": strings.TrimSpace(company),
		"ship_no": strings.TrimSpace(shipNo), "shipped_at": &now,
	}).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *PointsOrderRepository) Complete(id uint64) (*model.PointsExchangeOrder, error) {
	o, err := r.GetByID(id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if o.Status != model.PointsOrderPending && o.Status != model.PointsOrderShipped {
		return nil, errors.New("当前状态不可完成")
	}
	now := common.LocalTime(time.Now())
	if err := r.db.Model(o).Updates(map[string]interface{}{
		"status": model.PointsOrderCompleted, "completed_at": &now,
	}).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// CancelLocal 取消订单并退库存（积分退回由调用方走 user-service）
func (r *PointsOrderRepository) CancelLocal(id uint64, remark string) (*model.PointsExchangeOrder, error) {
	var out *model.PointsExchangeOrder
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var o model.PointsExchangeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&o, id).Error; err != nil {
			return errors.New("订单不存在")
		}
		if o.Status != model.PointsOrderPending {
			return errors.New("仅待发货订单可取消退积分")
		}
		now := common.LocalTime(time.Now())
		updates := map[string]interface{}{
			"status": model.PointsOrderCancelled, "cancelled_at": &now,
		}
		if strings.TrimSpace(remark) != "" {
			updates["admin_remark"] = strings.TrimSpace(remark)
		}
		if err := tx.Model(&o).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.PointsProduct{}).Where("id = ?", o.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", o.Quantity)).Error; err != nil {
			return err
		}
		out = &o
		out.Status = model.PointsOrderCancelled
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *PointsOrderRepository) Remark(id uint64, remark string) (*model.PointsExchangeOrder, error) {
	if _, err := r.GetByID(id); err != nil {
		return nil, errors.New("订单不存在")
	}
	if err := r.db.Model(&model.PointsExchangeOrder{}).Where("id = ?", id).
		Update("admin_remark", strings.TrimSpace(remark)).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *PointsOrderRepository) MapUserBriefs(ids []uint64) map[uint64][2]string {
	out := map[uint64][2]string{}
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
		out[u.ID] = [2]string{u.Nickname, u.Mobile}
	}
	return out
}
