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

func (r *MerchantRepository) ExpireDueThemeOrders(ctx context.Context) {
	now := time.Now()
	_ = r.db.WithContext(ctx).Model(&model.HomepageThemeOrder{}).
		Where("status = ? AND end_at < ?", model.ThemeOrderActive, now).
		Update("status", model.ThemeOrderExpired)
}

func (r *MerchantRepository) ListThemeSlots(ctx context.Context, onlyOn bool) ([]model.HomepageThemeSlot, error) {
	r.ExpireDueThemeOrders(ctx)
	q := r.db.WithContext(ctx).Model(&model.HomepageThemeSlot{})
	if onlyOn {
		q = q.Where("status = ?", model.ThemeSlotOn)
	}
	var list []model.HomepageThemeSlot
	err := q.Order("position ASC, sort ASC, id ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for i := range list {
		var o model.HomepageThemeOrder
		// 含尚未生效的排队单：取该坑 active 中最晚 end_at
		e := r.db.WithContext(ctx).Where("theme_slot_id = ? AND status = ? AND end_at > ?",
			list[i].ID, model.ThemeOrderActive, now).
			Order("end_at DESC").First(&o).Error
		if e == nil {
			list[i].HasActive = true
			list[i].OccupiedUntil = time.Time(o.EndAt).Format("2006-01-02 15:04:05")
		}
	}
	return list, nil
}

func (r *MerchantRepository) GetThemeSlot(ctx context.Context, id uint64) (*model.HomepageThemeSlot, error) {
	var s model.HomepageThemeSlot
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) UpdateThemeSlot(ctx context.Context, id uint64, updates map[string]interface{}) error {
	res := r.db.WithContext(ctx).Model(&model.HomepageThemeSlot{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("坑位不存在")
	}
	return nil
}

func (r *MerchantRepository) ListThemePackages(ctx context.Context, themeSlotID uint64, onlyOn bool) ([]model.HomepageThemePackage, error) {
	q := r.db.WithContext(ctx).Model(&model.HomepageThemePackage{})
	if onlyOn {
		q = q.Where("status = ?", model.ThemeSlotOn)
	}
	if themeSlotID > 0 {
		q = q.Where("theme_slot_id = 0 OR theme_slot_id = ?", themeSlotID)
	}
	var list []model.HomepageThemePackage
	err := q.Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *MerchantRepository) GetThemePackage(ctx context.Context, id uint64) (*model.HomepageThemePackage, error) {
	var p model.HomepageThemePackage
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MerchantRepository) CreateThemePackage(ctx context.Context, p *model.HomepageThemePackage) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *MerchantRepository) UpdateThemePackage(ctx context.Context, id uint64, updates map[string]interface{}) error {
	res := r.db.WithContext(ctx).Model(&model.HomepageThemePackage{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("套餐不存在")
	}
	return nil
}

func (r *MerchantRepository) ListThemeOrders(ctx context.Context, shopID, themeSlotID uint64, page, pageSize int) ([]model.HomepageThemeOrder, int64, error) {
	r.ExpireDueThemeOrders(ctx)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.HomepageThemeOrder{})
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	if themeSlotID > 0 {
		q = q.Where("theme_slot_id = ?", themeSlotID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.HomepageThemeOrder
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) ActiveThemeOrderForSlot(ctx context.Context, themeSlotID uint64) (*model.HomepageThemeOrder, error) {
	r.ExpireDueThemeOrders(ctx)
	now := time.Now()
	var o model.HomepageThemeOrder
	err := r.db.WithContext(ctx).Where("theme_slot_id = ? AND status = ? AND start_at <= ? AND end_at > ?",
		themeSlotID, model.ThemeOrderActive, now, now).
		Order("end_at DESC").First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *MerchantRepository) LatestThemeOrderQueueEnd(ctx context.Context, themeSlotID uint64) time.Time {
	r.ExpireDueThemeOrders(ctx)
	now := time.Now()
	var o model.HomepageThemeOrder
	err := r.db.WithContext(ctx).Where("theme_slot_id = ? AND status = ?", themeSlotID, model.ThemeOrderActive).
		Order("end_at DESC").First(&o).Error
	if err != nil {
		return now
	}
	end := time.Time(o.EndAt)
	if end.After(now) {
		return end
	}
	return now
}

// PurchaseThemeOrder 扣款或代开通；同坑位顺延排队
func (r *MerchantRepository) PurchaseThemeOrder(ctx context.Context, order *model.HomepageThemeOrder, deductWallet bool, operatorID *uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		_ = tx.Model(&model.HomepageThemeOrder{}).
			Where("status = ? AND end_at < ?", model.ThemeOrderActive, now).
			Update("status", model.ThemeOrderExpired)

		var prev model.HomepageThemeOrder
		err := tx.Where("theme_slot_id = ? AND status = ?", order.ThemeSlotID, model.ThemeOrderActive).
			Order("end_at DESC").First(&prev).Error
		start := now
		if err == nil && time.Time(prev.EndAt).After(now) {
			start = time.Time(prev.EndAt)
		}
		end := start.Add(time.Duration(order.DurationDays) * 24 * time.Hour)
		order.StartAt = common.LocalTime(start)
		order.EndAt = common.LocalTime(end)
		order.Status = model.ThemeOrderActive

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
				ChangeType:     model.WalletLogThemeSlot,
				Amount:         -order.Amount,
				BalanceAfter:   w.Balance,
				FrozenAfter:    w.FrozenBalance,
				DepositAfter:   w.Deposit,
				Remark:         "主题集市坑位购买",
				OperatorUserID: operatorID,
				RefType:        "homepage_theme_order",
				RefID:          order.ID,
			}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
			return tx.Model(order).Update("wallet_log_id", log.ID).Error
		}
		return tx.Create(order).Error
	})
}

func (r *MerchantRepository) CategoryExistsShow(ctx context.Context, id uint64) bool {
	var n int64
	r.db.WithContext(ctx).Table("product_categories").Where("id = ? AND is_show = 1", id).Count(&n)
	return n > 0
}

func (r *MerchantRepository) ProductOnSaleOfShop(ctx context.Context, productID, shopID uint64) bool {
	var n int64
	r.db.WithContext(ctx).Table("products").Where("id = ? AND shop_id = ? AND status = ?", productID, shopID, "on_sale").Count(&n)
	return n > 0
}

func (r *MerchantRepository) BuildThemeTiles(ctx context.Context) ([]model.ThemeTile, error) {
	r.ExpireDueThemeOrders(ctx)
	slots, err := r.ListThemeSlots(ctx, true)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	out := make([]model.ThemeTile, 0, len(slots))
	for _, s := range slots {
		tile := model.ThemeTile{
			Position: s.Position,
			SlotID:   s.ID,
			SlotKey:  s.SlotKey,
			Name:     s.Name,
			Desc:     s.Desc,
			CoverURL: s.CoverURL,
			LinkType: s.DefaultLinkType,
			LinkID:   s.DefaultLinkID,
			Paid:     false,
		}
		var o model.HomepageThemeOrder
		e := r.db.WithContext(ctx).Where("theme_slot_id = ? AND status = ? AND start_at <= ? AND end_at > ?",
			s.ID, model.ThemeOrderActive, now, now).
			Order("end_at DESC").First(&o).Error
		if e == nil {
			tile.Name = o.Title
			tile.Desc = o.Subtitle
			if o.CoverURL != "" {
				tile.CoverURL = o.CoverURL
			}
			tile.LinkType = o.LinkType
			tile.LinkID = o.LinkID
			tile.Paid = true
			tile.ShopID = o.ShopID
		}
		out = append(out, tile)
	}
	return out, nil
}
