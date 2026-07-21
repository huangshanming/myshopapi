package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *MerchantRepository) CreateCoupon(ctx context.Context, c *model.Coupon, scopes []model.CouponScope) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		for i := range scopes {
			scopes[i].CouponID = c.ID
			scopes[i].ID = 0
		}
		if len(scopes) > 0 {
			if err := tx.Create(&scopes).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *MerchantRepository) UpdateCoupon(ctx context.Context, id uint64, updates map[string]interface{}, scopes *[]model.CouponScope) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Coupon{}).Where("id = ?", id).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("优惠券不存在")
		}
		if scopes != nil {
			if err := tx.Where("coupon_id = ?", id).Delete(&model.CouponScope{}).Error; err != nil {
				return err
			}
			for i := range *scopes {
				(*scopes)[i].CouponID = id
				(*scopes)[i].ID = 0
			}
			if len(*scopes) > 0 {
				if err := tx.Create(scopes).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *MerchantRepository) GetCoupon(ctx context.Context, id uint64) (*model.Coupon, error) {
	var c model.Coupon
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	scopes, _ := r.ListCouponScopes(ctx, id)
	c.Scopes = scopes
	return &c, nil
}

func (r *MerchantRepository) ListCouponScopes(ctx context.Context, couponID uint64) ([]model.CouponScope, error) {
	var list []model.CouponScope
	err := r.db.WithContext(ctx).Where("coupon_id = ?", couponID).Find(&list).Error
	return list, err
}

func (r *MerchantRepository) ListCoupons(ctx context.Context, issuerType string, shopID uint64, status, keyword string, page, pageSize int) ([]model.Coupon, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.Coupon{})
	if issuerType != "" {
		q = q.Where("issuer_type = ?", issuerType)
	}
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	} else if issuerType == model.CouponIssuerPlatform {
		q = q.Where("shop_id = 0")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Coupon
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	now := time.Now()
	for i := range list {
		scopes, _ := r.ListCouponScopes(ctx, list[i].ID)
		list[i].Scopes = scopes
		list[i].DisplayStatus = couponDisplayStatus(&list[i], now)
		if list[i].TotalCount > 0 {
			rem := list[i].TotalCount - list[i].ClaimedCount
			if rem < 0 {
				rem = 0
			}
			list[i].Remaining = &rem
		}
	}
	return list, total, nil
}

func couponDisplayStatus(c *model.Coupon, now time.Time) string {
	if c.Status == model.CouponStatusOff {
		return "off"
	}
	if c.Status == model.CouponStatusDraft {
		return "draft"
	}
	if c.Status == model.CouponStatusExpired {
		return "expired"
	}
	if c.ValidType == model.CouponValidFixed && c.ValidEnd != nil && time.Time(*c.ValidEnd).Before(now) {
		return "expired"
	}
	if c.ValidType == model.CouponValidFixed && c.ValidStart != nil && time.Time(*c.ValidStart).After(now) {
		return "not_started"
	}
	if c.TotalCount > 0 && c.ClaimedCount >= c.TotalCount {
		return "sold_out"
	}
	if c.Status == model.CouponStatusOn {
		return "active"
	}
	return c.Status
}

func (r *MerchantRepository) CountUserClaims(ctx context.Context, couponID, userID uint64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.UserCoupon{}).Where("coupon_id = ? AND user_id = ?", couponID, userID).Count(&n).Error
	return n, err
}

func (r *MerchantRepository) UserCreatedAt(ctx context.Context, userID uint64) (time.Time, error) {
	var createdAt time.Time
	err := r.db.WithContext(ctx).Table("users").Select("created_at").Where("id = ?", userID).Scan(&createdAt).Error
	return createdAt, err
}

type productLite struct {
	ID         uint64
	ShopID     uint64
	CategoryID uint64
	Status     string
}

func (r *MerchantRepository) GetProductsLite(ctx context.Context, ids []uint64) (map[uint64]productLite, error) {
	out := map[uint64]productLite{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []struct {
		ID         uint64 `gorm:"column:id"`
		ShopID     uint64 `gorm:"column:shop_id"`
		CategoryID uint64 `gorm:"column:category_id"`
		Status     string `gorm:"column:status"`
	}
	if err := r.db.WithContext(ctx).Table("products").Select("id, shop_id, category_id, status").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = productLite{ID: row.ID, ShopID: row.ShopID, CategoryID: row.CategoryID, Status: row.Status}
	}
	return out, nil
}

func (r *MerchantRepository) ClaimCoupon(ctx context.Context, userID uint64, c *model.Coupon, source, batchNo string) (*model.UserCoupon, error) {
	var uc model.UserCoupon
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var locked model.Coupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&locked, c.ID).Error; err != nil {
			return errors.New("优惠券不存在")
		}
		if locked.Status != model.CouponStatusOn {
			return errors.New("优惠券不可领取")
		}
		if locked.TotalCount > 0 && locked.ClaimedCount >= locked.TotalCount {
			return errors.New("已领完")
		}
		var claimed int64
		if err := tx.Model(&model.UserCoupon{}).Where("coupon_id = ? AND user_id = ?", c.ID, userID).Count(&claimed).Error; err != nil {
			return err
		}
		if int(claimed) >= locked.PerUserLimit {
			return errors.New("已达限领")
		}
		now := time.Now()
		start, end, err := calcUserCouponValid(&locked, now)
		if err != nil {
			return err
		}
		uc = model.UserCoupon{
			CouponID:     locked.ID,
			UserID:       userID,
			ShopID:       locked.ShopID,
			Status:       model.UserCouponUnused,
			Source:       source,
			ValidStart:   common.LocalTime(start),
			ValidEnd:     common.LocalTime(end),
			ClaimBatchNo: batchNo,
		}
		if err := tx.Create(&uc).Error; err != nil {
			return err
		}
		res := tx.Model(&model.Coupon{}).Where("id = ? AND (total_count = 0 OR claimed_count < total_count)", locked.ID).
			UpdateColumn("claimed_count", gorm.Expr("claimed_count + 1"))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("已领完")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &uc, nil
}

func calcUserCouponValid(c *model.Coupon, now time.Time) (time.Time, time.Time, error) {
	if c.ValidType == model.CouponValidRelative {
		if c.ValidDays < 1 {
			return time.Time{}, time.Time{}, errors.New("有效天数无效")
		}
		start := now
		end := now.Add(time.Duration(c.ValidDays) * 24 * time.Hour)
		return start, end, nil
	}
	if c.ValidStart == nil || c.ValidEnd == nil {
		return time.Time{}, time.Time{}, errors.New("固定有效期未配置")
	}
	start := time.Time(*c.ValidStart)
	end := time.Time(*c.ValidEnd)
	if end.Before(now) {
		return time.Time{}, time.Time{}, errors.New("优惠券已过期")
	}
	if start.After(now) {
		// 领取后从配置开始日起算
	} else {
		start = now
	}
	return start, end, nil
}

func (r *MerchantRepository) ListUserCoupons(ctx context.Context, userID uint64, status string, page, pageSize int) ([]model.UserCoupon, int64, error) {
	r.ExpireUserCoupons(ctx, userID)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.UserCoupon{}).Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.UserCoupon
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		if c, e := r.GetCoupon(ctx, list[i].CouponID); e == nil {
			list[i].Coupon = c
			list[i].CouponName = c.Name
		}
	}
	return list, total, nil
}

func (r *MerchantRepository) ExpireUserCoupons(ctx context.Context, userID uint64) {
	now := time.Now()
	q := r.db.WithContext(ctx).Model(&model.UserCoupon{}).
		Where("status = ? AND valid_end < ?", model.UserCouponUnused, now)
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	_ = q.Update("status", model.UserCouponExpired)
}

func (r *MerchantRepository) GetUserCoupon(ctx context.Context, id uint64) (*model.UserCoupon, error) {
	var uc model.UserCoupon
	if err := r.db.WithContext(ctx).First(&uc, id).Error; err != nil {
		return nil, err
	}
	if c, e := r.GetCoupon(ctx, uc.CouponID); e == nil {
		uc.Coupon = c
		uc.CouponName = c.Name
	}
	return &uc, nil
}

func (r *MerchantRepository) LockUserCoupon(ctx context.Context, userCouponID, userID, orderID uint64, discount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var uc model.UserCoupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&uc, userCouponID).Error; err != nil {
			return errors.New("优惠券不存在")
		}
		if uc.UserID != userID {
			return errors.New("优惠券不属于当前用户")
		}
		now := time.Now()
		if uc.Status != model.UserCouponUnused {
			return errors.New("优惠券不可用")
		}
		if time.Time(uc.ValidEnd).Before(now) {
			_ = tx.Model(&uc).Update("status", model.UserCouponExpired)
			return errors.New("优惠券已过期")
		}
		lt := common.LocalTime(now)
		return tx.Model(&uc).Updates(map[string]interface{}{
			"status":          model.UserCouponLocked,
			"order_id":        orderID,
			"locked_at":       lt,
			"discount_amount": discount,
		}).Error
	})
}

func (r *MerchantRepository) UnlockUserCoupon(ctx context.Context, userCouponID, orderID uint64, action string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var uc model.UserCoupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&uc, userCouponID).Error; err != nil {
			return errors.New("优惠券不存在")
		}
		if uc.Status != model.UserCouponLocked && uc.Status != model.UserCouponUsed {
			return nil
		}
		if orderID > 0 && uc.OrderID != orderID {
			return errors.New("订单与优惠券不匹配")
		}
		now := time.Now()
		status := model.UserCouponUnused
		if time.Time(uc.ValidEnd).Before(now) {
			status = model.UserCouponExpired
		}
		if err := tx.Model(&uc).Updates(map[string]interface{}{
			"status":          status,
			"order_id":        0,
			"locked_at":       nil,
			"used_at":         nil,
			"discount_amount": 0,
		}).Error; err != nil {
			return err
		}
		log := model.CouponRedeemLog{
			UserCouponID:   uc.ID,
			CouponID:       uc.CouponID,
			UserID:         uc.UserID,
			OrderID:        orderID,
			ShopID:         uc.ShopID,
			DiscountAmount: uc.DiscountAmount,
			Action:         action,
		}
		return tx.Create(&log).Error
	})
}

func (r *MerchantRepository) RedeemUserCoupon(ctx context.Context, userCouponID, orderID uint64, discount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var uc model.UserCoupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&uc, userCouponID).Error; err != nil {
			return errors.New("优惠券不存在")
		}
		if uc.Status != model.UserCouponLocked {
			return errors.New("优惠券未锁定，无法核销")
		}
		if orderID > 0 && uc.OrderID != orderID {
			return errors.New("订单与优惠券不匹配")
		}
		now := common.LocalTime(time.Now())
		if err := tx.Model(&uc).Updates(map[string]interface{}{
			"status":          model.UserCouponUsed,
			"used_at":         now,
			"discount_amount": discount,
		}).Error; err != nil {
			return err
		}
		log := model.CouponRedeemLog{
			UserCouponID:   uc.ID,
			CouponID:       uc.CouponID,
			UserID:         uc.UserID,
			OrderID:        orderID,
			ShopID:         uc.ShopID,
			DiscountAmount: discount,
			Action:         model.CouponActionRedeem,
		}
		return tx.Create(&log).Error
	})
}

func (r *MerchantRepository) ListClaims(ctx context.Context, couponID uint64, page, pageSize int) ([]model.UserCoupon, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.UserCoupon{}).Where("coupon_id = ?", couponID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.UserCoupon
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) ListRedeems(ctx context.Context, couponID uint64, page, pageSize int) ([]model.CouponRedeemLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.CouponRedeemLog{}).Where("coupon_id = ? AND action = ?", couponID, model.CouponActionRedeem)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CouponRedeemLog
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) CouponStats(ctx context.Context, couponID uint64) (map[string]interface{}, error) {
	var claimed int64
	_ = r.db.Model(&model.UserCoupon{}).Where("coupon_id = ?", couponID).Count(&claimed)
	var redeemed int64
	_ = r.db.Model(&model.UserCoupon{}).Where("coupon_id = ? AND status = ?", couponID, model.UserCouponUsed).Count(&redeemed)
	var sum float64
	_ = r.db.Model(&model.CouponRedeemLog{}).
		Select("COALESCE(SUM(discount_amount),0)").
		Where("coupon_id = ? AND action = ?", couponID, model.CouponActionRedeem).
		Scan(&sum)
	rate := 0.0
	if claimed > 0 {
		rate = float64(redeemed) / float64(claimed) * 100
	}
	return map[string]interface{}{
		"claimed_count":  claimed,
		"redeemed_count": redeemed,
		"redeem_rate":    fmt.Sprintf("%.2f", rate),
		"discount_total": sum,
	}, nil
}

func (r *MerchantRepository) CreateGrant(ctx context.Context, g *model.CouponGrant) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *MerchantRepository) ListCenterCoupons(ctx context.Context, shopID uint64) ([]model.Coupon, error) {
	now := time.Now()
	var list []model.Coupon
	q := r.db.WithContext(ctx).Where("status = ?", model.CouponStatusOn)
	// 平台券 + 可选店铺券
	if shopID > 0 {
		q = q.Where("(issuer_type = ? AND shop_id = 0) OR (issuer_type = ? AND shop_id = ?)",
			model.CouponIssuerPlatform, model.CouponIssuerShop, shopID)
	} else {
		q = q.Where("issuer_type = ? AND shop_id = 0", model.CouponIssuerPlatform)
	}
	err := q.Order("id DESC").Limit(100).Find(&list).Error
	if err != nil {
		return nil, err
	}
	out := make([]model.Coupon, 0, len(list))
	for i := range list {
		ds := couponDisplayStatus(&list[i], now)
		if ds != "active" && ds != "sold_out" {
			continue
		}
		// 仅展示含 direct/popup 渠道
		has := false
		for _, ch := range list[i].Channels {
			if ch == model.CouponChannelDirect || ch == model.CouponChannelPopup {
				has = true
				break
			}
		}
		if !has && len(list[i].Channels) == 0 {
			has = true
		}
		if !has {
			continue
		}
		list[i].DisplayStatus = ds
		scopes, _ := r.ListCouponScopes(ctx, list[i].ID)
		list[i].Scopes = scopes
		if list[i].TotalCount > 0 {
			rem := list[i].TotalCount - list[i].ClaimedCount
			if rem < 0 {
				rem = 0
			}
			list[i].Remaining = &rem
		}
		out = append(out, list[i])
	}
	return out, nil
}

func (r *MerchantRepository) ListPopupCoupons(ctx context.Context) ([]model.Coupon, error) {
	list, err := r.ListCenterCoupons(ctx, 0)
	if err != nil {
		return nil, err
	}
	out := make([]model.Coupon, 0)
	for _, c := range list {
		for _, ch := range c.Channels {
			if ch == model.CouponChannelPopup {
				out = append(out, c)
				break
			}
		}
	}
	return out, nil
}

func (r *MerchantRepository) ListUserUnusedCoupons(ctx context.Context, userID uint64) ([]model.UserCoupon, error) {
	r.ExpireUserCoupons(ctx, userID)
	var list []model.UserCoupon
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, model.UserCouponUnused).
		Order("valid_end ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	for i := range list {
		if c, e := r.GetCoupon(ctx, list[i].CouponID); e == nil {
			list[i].Coupon = c
			list[i].CouponName = c.Name
		}
	}
	return list, nil
}

func (r *MerchantRepository) ListOrderGiftCoupons(ctx context.Context, shopID uint64) ([]model.Coupon, error) {
	now := time.Now()
	var list []model.Coupon
	err := r.db.WithContext(ctx).Where("status = ?", model.CouponStatusOn).
		Where("(issuer_type = ? AND shop_id = 0) OR (issuer_type = ? AND shop_id = ?)",
			model.CouponIssuerPlatform, model.CouponIssuerShop, shopID).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	out := make([]model.Coupon, 0)
	for i := range list {
		if couponDisplayStatus(&list[i], now) != "active" {
			continue
		}
		for _, ch := range list[i].Channels {
			if ch == model.CouponChannelOrderGift {
				scopes, _ := r.ListCouponScopes(ctx, list[i].ID)
				list[i].Scopes = scopes
				out = append(out, list[i])
				break
			}
		}
	}
	return out, nil
}
