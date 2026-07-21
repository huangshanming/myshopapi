package repository

import (
	"context"
	"errors"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *MerchantRepository) ExpireDueSlotOrders(ctx context.Context) {
	now := time.Now()
	_ = r.db.WithContext(ctx).Model(&model.HomepageSlotOrder{}).
		Where("status = ? AND end_at < ?", model.SlotOrderActive, now).
		Update("status", model.SlotOrderExpired)
}

func (r *MerchantRepository) ListSlotPackages(ctx context.Context, slotType string, onlyOn bool) ([]model.HomepageSlotPackage, error) {
	q := r.db.WithContext(ctx).Model(&model.HomepageSlotPackage{})
	if slotType != "" {
		q = q.Where("slot_type = ?", slotType)
	}
	if onlyOn {
		q = q.Where("status = ?", model.SlotPkgOn)
	}
	var list []model.HomepageSlotPackage
	err := q.Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *MerchantRepository) GetSlotPackage(ctx context.Context, id uint64) (*model.HomepageSlotPackage, error) {
	var p model.HomepageSlotPackage
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MerchantRepository) CreateSlotPackage(ctx context.Context, p *model.HomepageSlotPackage) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *MerchantRepository) UpdateSlotPackage(ctx context.Context, p *model.HomepageSlotPackage) error {
	return r.db.WithContext(ctx).Model(p).Updates(map[string]interface{}{
		"slot_type": p.SlotType, "name": p.Name, "price": p.Price,
		"duration_days": p.DurationDays, "status": p.Status, "sort": p.Sort, "remark": p.Remark,
	}).Error
}

func (r *MerchantRepository) ListSlotSettings(ctx context.Context) ([]model.HomepageSlotSetting, error) {
	var list []model.HomepageSlotSetting
	err := r.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *MerchantRepository) GetSlotSetting(ctx context.Context, slotType string) (*model.HomepageSlotSetting, error) {
	var s model.HomepageSlotSetting
	err := r.db.WithContext(ctx).Where("slot_type = ?", slotType).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s = model.HomepageSlotSetting{SlotType: slotType, HomeLimit: 8}
		_ = r.db.WithContext(ctx).Create(&s).Error
		return &s, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) UpsertSlotSetting(ctx context.Context, slotType string, homeLimit int) error {
	var s model.HomepageSlotSetting
	err := r.db.WithContext(ctx).Where("slot_type = ?", slotType).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.WithContext(ctx).Create(&model.HomepageSlotSetting{SlotType: slotType, HomeLimit: homeLimit}).Error
	}
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&s).Update("home_limit", homeLimit).Error
}

func (r *MerchantRepository) ListSlotOrders(ctx context.Context, shopID uint64, slotType, status string, page, pageSize int) ([]model.HomepageSlotOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.HomepageSlotOrder{})
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	if slotType != "" {
		q = q.Where("slot_type = ?", slotType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.HomepageSlotOrder
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) LatestActiveSlotOrder(ctx context.Context, shopID uint64, slotType string, targetID uint64) (*model.HomepageSlotOrder, error) {
	r.ExpireDueSlotOrders(ctx)
	q := r.db.WithContext(ctx).Where("shop_id = ? AND slot_type = ? AND status = ?", shopID, slotType, model.SlotOrderActive)
	if slotType == model.SlotArticle {
		q = q.Where("target_id = ?", targetID)
	} else {
		q = q.Where("target_id = ?", shopID)
	}
	var o model.HomepageSlotOrder
	err := q.Order("end_at DESC").First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// PurchaseSlotOrder 钱包扣款或超管开通，立即生效（可顺延）
func (r *MerchantRepository) PurchaseSlotOrder(ctx context.Context, order *model.HomepageSlotOrder, deductWallet bool, operatorID *uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		_ = tx.Model(&model.HomepageSlotOrder{}).
			Where("status = ? AND end_at < ?", model.SlotOrderActive, now).
			Update("status", model.SlotOrderExpired)

		targetID := order.TargetID
		q := tx.Where("shop_id = ? AND slot_type = ? AND status = ?", order.ShopID, order.SlotType, model.SlotOrderActive)
		if order.SlotType == model.SlotArticle {
			q = q.Where("target_id = ?", targetID)
		} else {
			q = q.Where("target_id = ?", order.ShopID)
		}
		var prev model.HomepageSlotOrder
		err := q.Order("end_at DESC").First(&prev).Error
		start := now
		if err == nil && time.Time(prev.EndAt).After(now) {
			start = time.Time(prev.EndAt)
		}
		end := start.Add(time.Duration(order.DurationDays) * 24 * time.Hour)
		order.StartAt = common.LocalTime(start)
		order.EndAt = common.LocalTime(end)
		order.Status = model.SlotOrderActive
		if order.SlotType != model.SlotArticle {
			order.TargetID = order.ShopID
		}

		var walletLogID uint64
		if deductWallet {
			var w model.ShopWallet
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("shop_id = ?", order.ShopID).First(&w).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w = model.ShopWallet{ShopID: order.ShopID}
				if err := tx.Create(&w).Error; err != nil {
					return err
				}
				err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("shop_id = ?", order.ShopID).First(&w).Error
			}
			if err != nil {
				return err
			}
			if w.Balance+0.0001 < order.Amount {
				return errors.New("余额不足，请联系平台充值后再购买")
			}
			w.Balance -= order.Amount
			if err := tx.Model(&w).Update("balance", w.Balance).Error; err != nil {
				return err
			}
			if err := tx.Create(order).Error; err != nil {
				return err
			}
			log := model.ShopWalletLog{
				ShopID:         order.ShopID,
				ChangeType:     model.WalletLogHomepageSlot,
				Amount:         -order.Amount,
				BalanceAfter:   w.Balance,
				FrozenAfter:    w.FrozenBalance,
				DepositAfter:   w.Deposit,
				Remark:         "首页展位套餐购买",
				OperatorUserID: operatorID,
				RefType:        "homepage_slot_order",
				RefID:          order.ID,
			}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
			walletLogID = log.ID
			return tx.Model(order).Update("wallet_log_id", walletLogID).Error
		}
		return tx.Create(order).Error
	})
}

func (r *MerchantRepository) ListPublicShopsWithSlot(ctx context.Context, slotType string, page, pageSize, homeLimit int) ([]model.Shop, int64, error) {
	r.ExpireDueSlotOrders(ctx)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	now := time.Now()

	type row struct {
		model.Shop
		Boost int `gorm:"column:boost"`
	}

	baseFrom := `
FROM shops s
LEFT JOIN (
  SELECT DISTINCT target_id
  FROM homepage_slot_orders
  WHERE slot_type = ? AND status = ? AND start_at <= ? AND end_at > ?
) o ON o.target_id = s.id
WHERE s.status = ?`

	countSQL := "SELECT COUNT(*) " + baseFrom
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, slotType, model.SlotOrderActive, now, now, model.ShopApproved).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := pageSize
	offset := (page - 1) * pageSize
	if homeLimit > 0 {
		limit = homeLimit
		offset = 0
		if int64(homeLimit) < total {
			total = int64(homeLimit)
		}
	}

	listSQL := `
SELECT s.*, CASE WHEN o.target_id IS NULL THEN 0 ELSE 1 END AS boost
` + baseFrom + `
ORDER BY boost DESC, s.id DESC
LIMIT ? OFFSET ?`

	var rows []row
	if err := r.db.WithContext(ctx).Raw(listSQL, slotType, model.SlotOrderActive, now, now, model.ShopApproved, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]model.Shop, 0, len(rows))
	for _, x := range rows {
		out = append(out, x.Shop)
	}
	return out, total, nil
}

// GetArticleTitle 同库查文章标题（catalog 表 community_article）
func (r *MerchantRepository) GetArticleTitle(ctx context.Context, id uint64) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("community_article").Select("title").Where("id = ?", id).Scan(&title).Error
	if err != nil {
		return "", err
	}
	if title == "" {
		return "", gorm.ErrRecordNotFound
	}
	return title, nil
}

// ActivePaidTargetIDs 返回当前生效的 target_id 集合（用于标记 paid）
func (r *MerchantRepository) ActivePaidTargetIDs(ctx context.Context, slotType string) (map[uint64]bool, error) {
	r.ExpireDueSlotOrders(ctx)
	now := time.Now()
	var ids []uint64
	err := r.db.WithContext(ctx).Model(&model.HomepageSlotOrder{}).
		Where("slot_type = ? AND status = ? AND start_at <= ? AND end_at > ?", slotType, model.SlotOrderActive, now, now).
		Distinct("target_id").Pluck("target_id", &ids).Error
	m := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m, err
}
