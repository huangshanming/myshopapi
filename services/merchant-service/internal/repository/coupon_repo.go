package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	couponColumns      = "id, IFNULL(name,'') AS name, issuer_type, shop_id, coupon_type, threshold_amount, discount_amount, discount_rate, max_discount_amount, scope_type, total_count, claimed_count, per_user_limit, valid_type, valid_start, valid_end, valid_days, stackable, user_identity, IFNULL(channels,'') AS channels, status, IFNULL(remark,'') AS remark, created_by, created_at, updated_at"
	couponScopeColumns = "id, coupon_id, ref_type, ref_id"
	userCouponColumns  = "id, coupon_id, user_id, shop_id, status, source, valid_start, valid_end, order_id, locked_at, used_at, IFNULL(claim_batch_no,'') AS claim_batch_no, discount_amount, created_at, updated_at"
)

func (r *MerchantRepository) CreateCoupon(ctx context.Context, c *model.Coupon, scopes []model.CouponScope) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		res, err := session.ExecCtx(ctx,
			`INSERT INTO coupons (name, issuer_type, shop_id, coupon_type, threshold_amount, discount_amount, discount_rate, max_discount_amount, scope_type, total_count, claimed_count, per_user_limit, valid_type, valid_start, valid_end, valid_days, stackable, user_identity, channels, status, remark, created_by)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			c.Name, c.IssuerType, c.ShopID, c.CouponType, c.ThresholdAmount, c.DiscountAmount, c.DiscountRate,
			c.MaxDiscountAmount, c.ScopeType, c.TotalCount, c.ClaimedCount, c.PerUserLimit, c.ValidType,
			c.ValidStart, c.ValidEnd, c.ValidDays, c.Stackable, c.UserIdentity, c.Channels, c.Status, c.Remark, c.CreatedBy,
		)
		if err != nil {
			return err
		}
		id, err := lastInsertID(res)
		if err != nil {
			return err
		}
		c.ID = id
		for i := range scopes {
			scopes[i].CouponID = c.ID
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO coupon_scopes (coupon_id, ref_type, ref_id) VALUES (?,?,?)",
				scopes[i].CouponID, scopes[i].RefType, scopes[i].RefID,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *MerchantRepository) UpdateCoupon(ctx context.Context, id uint64, updates map[string]interface{}, scopes *[]model.CouponScope) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		query, args, err := buildUpdate("coupons", updates, "id=?", id)
		if err != nil {
			return err
		}
		n, err := execRows(ctx, session, query, args...)
		if err != nil {
			return err
		}
		if n == 0 {
			return errors.New("优惠券不存在")
		}
		if scopes != nil {
			if _, err := session.ExecCtx(ctx, "DELETE FROM coupon_scopes WHERE coupon_id=?", id); err != nil {
				return err
			}
			for i := range *scopes {
				(*scopes)[i].CouponID = id
				if _, err := session.ExecCtx(ctx,
					"INSERT INTO coupon_scopes (coupon_id, ref_type, ref_id) VALUES (?,?,?)",
					(*scopes)[i].CouponID, (*scopes)[i].RefType, (*scopes)[i].RefID,
				); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *MerchantRepository) GetCoupon(ctx context.Context, id uint64) (*model.Coupon, error) {
	var c model.Coupon
	if err := r.conn.QueryRowPartialCtx(ctx, &c,
		"SELECT "+couponColumns+" FROM coupons WHERE id=? LIMIT 1", id,
	); err != nil {
		return nil, err
	}
	scopes, _ := r.ListCouponScopes(ctx, id)
	c.Scopes = scopes
	return &c, nil
}

func (r *MerchantRepository) ListCouponScopes(ctx context.Context, couponID uint64) ([]model.CouponScope, error) {
	var list []model.CouponScope
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+couponScopeColumns+" FROM coupon_scopes WHERE coupon_id=?", couponID,
	)
	return list, err
}

func (r *MerchantRepository) ListCoupons(ctx context.Context, issuerType string, shopID uint64, status, keyword string, page, pageSize int) ([]model.Coupon, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	where := "1=1"
	args := make([]any, 0, 4)
	if issuerType != "" {
		where += " AND issuer_type=?"
		args = append(args, issuerType)
	}
	if shopID > 0 {
		where += " AND shop_id=?"
		args = append(args, shopID)
	} else if issuerType == model.CouponIssuerPlatform {
		where += " AND shop_id=0"
	}
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}
	if keyword != "" {
		where += " AND name LIKE ?"
		args = append(args, "%"+keyword+"%")
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM coupons WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.Coupon
	if err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+couponColumns+" FROM coupons WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	); err != nil {
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
	return countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_coupons WHERE coupon_id=? AND user_id=?", couponID, userID,
	)
}

func (r *MerchantRepository) UserCreatedAt(ctx context.Context, userID uint64) (time.Time, error) {
	var createdAt time.Time
	err := r.conn.QueryRowPartialCtx(ctx, &createdAt,
		"SELECT created_at FROM users WHERE id=? LIMIT 1", userID,
	)
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
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(
		"SELECT id, shop_id, category_id, status FROM products WHERE id IN (%s)",
		strings.Join(placeholders, ","),
	)
	var rows []struct {
		ID         uint64 `db:"id"`
		ShopID     uint64 `db:"shop_id"`
		CategoryID uint64 `db:"category_id"`
		Status     string `db:"status"`
	}
	if err := r.conn.QueryRowsPartialCtx(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = productLite{ID: row.ID, ShopID: row.ShopID, CategoryID: row.CategoryID, Status: row.Status}
	}
	return out, nil
}

func (r *MerchantRepository) ClaimCoupon(ctx context.Context, userID uint64, c *model.Coupon, source, batchNo string) (*model.UserCoupon, error) {
	var uc model.UserCoupon
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var locked model.Coupon
		if err := session.QueryRowPartialCtx(ctx, &locked,
			"SELECT "+couponColumns+" FROM coupons WHERE id=? FOR UPDATE", c.ID,
		); err != nil {
			return errors.New("优惠券不存在")
		}
		if locked.Status != model.CouponStatusOn {
			return errors.New("优惠券不可领取")
		}
		if locked.TotalCount > 0 && locked.ClaimedCount >= locked.TotalCount {
			return errors.New("已领完")
		}
		claimed, err := countQuery(ctx, session,
			"SELECT COUNT(*) FROM user_coupons WHERE coupon_id=? AND user_id=?", c.ID, userID,
		)
		if err != nil {
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
		res, err := session.ExecCtx(ctx,
			`INSERT INTO user_coupons (coupon_id, user_id, shop_id, status, source, valid_start, valid_end, claim_batch_no)
			 VALUES (?,?,?,?,?,?,?,?)`,
			uc.CouponID, uc.UserID, uc.ShopID, uc.Status, uc.Source, uc.ValidStart, uc.ValidEnd, uc.ClaimBatchNo,
		)
		if err != nil {
			return err
		}
		ucID, err := lastInsertID(res)
		if err != nil {
			return err
		}
		uc.ID = ucID
		n, err := execRows(ctx, session,
			"UPDATE coupons SET claimed_count=claimed_count+1 WHERE id=? AND (total_count=0 OR claimed_count<total_count)",
			locked.ID,
		)
		if err != nil {
			return err
		}
		if n == 0 {
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
	where := "user_id=?"
	args := []any{userID}
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM user_coupons WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.UserCoupon
	if err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+userCouponColumns+" FROM user_coupons WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	); err != nil {
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
	if userID > 0 {
		_, _ = r.conn.ExecCtx(ctx,
			"UPDATE user_coupons SET status=? WHERE status=? AND valid_end<? AND user_id=?",
			model.UserCouponExpired, model.UserCouponUnused, now, userID,
		)
		return
	}
	_, _ = r.conn.ExecCtx(ctx,
		"UPDATE user_coupons SET status=? WHERE status=? AND valid_end<?",
		model.UserCouponExpired, model.UserCouponUnused, now,
	)
}

func (r *MerchantRepository) GetUserCoupon(ctx context.Context, id uint64) (*model.UserCoupon, error) {
	var uc model.UserCoupon
	if err := r.conn.QueryRowPartialCtx(ctx, &uc,
		"SELECT "+userCouponColumns+" FROM user_coupons WHERE id=? LIMIT 1", id,
	); err != nil {
		return nil, err
	}
	if c, e := r.GetCoupon(ctx, uc.CouponID); e == nil {
		uc.Coupon = c
		uc.CouponName = c.Name
	}
	return &uc, nil
}

func (r *MerchantRepository) LockUserCoupon(ctx context.Context, userCouponID, userID, orderID uint64, discount float64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var uc model.UserCoupon
		if err := session.QueryRowPartialCtx(ctx, &uc,
			"SELECT "+userCouponColumns+" FROM user_coupons WHERE id=? FOR UPDATE", userCouponID,
		); err != nil {
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
			_, _ = session.ExecCtx(ctx, "UPDATE user_coupons SET status=? WHERE id=?", model.UserCouponExpired, uc.ID)
			return errors.New("优惠券已过期")
		}
		lt := common.LocalTime(now)
		_, err := session.ExecCtx(ctx,
			"UPDATE user_coupons SET status=?, order_id=?, locked_at=?, discount_amount=? WHERE id=?",
			model.UserCouponLocked, orderID, lt, discount, uc.ID,
		)
		return err
	})
}

func (r *MerchantRepository) UnlockUserCoupon(ctx context.Context, userCouponID, orderID uint64, action string) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var uc model.UserCoupon
		if err := session.QueryRowPartialCtx(ctx, &uc,
			"SELECT "+userCouponColumns+" FROM user_coupons WHERE id=? FOR UPDATE", userCouponID,
		); err != nil {
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
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_coupons SET status=?, order_id=0, locked_at=NULL, used_at=NULL, discount_amount=0 WHERE id=?",
			status, uc.ID,
		); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx,
			`INSERT INTO coupon_redeem_logs (user_coupon_id, coupon_id, user_id, order_id, shop_id, discount_amount, action)
			 VALUES (?,?,?,?,?,?,?)`,
			uc.ID, uc.CouponID, uc.UserID, orderID, uc.ShopID, uc.DiscountAmount, action,
		)
		return err
	})
}

func (r *MerchantRepository) RedeemUserCoupon(ctx context.Context, userCouponID, orderID uint64, discount float64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var uc model.UserCoupon
		if err := session.QueryRowPartialCtx(ctx, &uc,
			"SELECT "+userCouponColumns+" FROM user_coupons WHERE id=? FOR UPDATE", userCouponID,
		); err != nil {
			return errors.New("优惠券不存在")
		}
		if uc.Status != model.UserCouponLocked {
			return errors.New("优惠券未锁定，无法核销")
		}
		if orderID > 0 && uc.OrderID != orderID {
			return errors.New("订单与优惠券不匹配")
		}
		now := common.LocalTime(time.Now())
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_coupons SET status=?, used_at=?, discount_amount=? WHERE id=?",
			model.UserCouponUsed, now, discount, uc.ID,
		); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx,
			`INSERT INTO coupon_redeem_logs (user_coupon_id, coupon_id, user_id, order_id, shop_id, discount_amount, action)
			 VALUES (?,?,?,?,?,?,?)`,
			uc.ID, uc.CouponID, uc.UserID, orderID, uc.ShopID, discount, model.CouponActionRedeem,
		)
		return err
	})
}

func (r *MerchantRepository) ListClaims(ctx context.Context, couponID uint64, page, pageSize int) ([]model.UserCoupon, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_coupons WHERE coupon_id=?", couponID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.UserCoupon
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+userCouponColumns+" FROM user_coupons WHERE coupon_id=? ORDER BY id DESC LIMIT ? OFFSET ?",
		couponID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) ListRedeems(ctx context.Context, couponID uint64, page, pageSize int) ([]model.CouponRedeemLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM coupon_redeem_logs WHERE coupon_id=? AND action=?",
		couponID, model.CouponActionRedeem,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.CouponRedeemLog
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		`SELECT id, user_coupon_id, coupon_id, user_id, order_id, shop_id, discount_amount, action, created_at
		 FROM coupon_redeem_logs WHERE coupon_id=? AND action=? ORDER BY id DESC LIMIT ? OFFSET ?`,
		couponID, model.CouponActionRedeem, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) CouponStats(ctx context.Context, couponID uint64) (map[string]interface{}, error) {
	claimed, _ := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_coupons WHERE coupon_id=?", couponID,
	)
	redeemed, _ := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_coupons WHERE coupon_id=? AND status=?", couponID, model.UserCouponUsed,
	)
	var sum float64
	_ = r.conn.QueryRowPartialCtx(ctx, &sum,
		"SELECT COALESCE(SUM(discount_amount),0) FROM coupon_redeem_logs WHERE coupon_id=? AND action=?",
		couponID, model.CouponActionRedeem,
	)
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
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO coupon_grants (coupon_id, operator_id, issuer_type, shop_id, user_count, success_count, batch_no)
		 VALUES (?,?,?,?,?,?,?)`,
		g.CouponID, g.OperatorID, g.IssuerType, g.ShopID, g.UserCount, g.SuccessCount, g.BatchNo,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	g.ID = id
	return nil
}

func (r *MerchantRepository) ListCenterCoupons(ctx context.Context, shopID uint64) ([]model.Coupon, error) {
	now := time.Now()
	where := "status=?"
	args := []any{model.CouponStatusOn}
	if shopID > 0 {
		where += " AND ((issuer_type=? AND shop_id=0) OR (issuer_type=? AND shop_id=?))"
		args = append(args, model.CouponIssuerPlatform, model.CouponIssuerShop, shopID)
	} else {
		where += " AND issuer_type=? AND shop_id=0"
		args = append(args, model.CouponIssuerPlatform)
	}
	var list []model.Coupon
	if err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+couponColumns+" FROM coupons WHERE "+where+" ORDER BY id DESC LIMIT 100",
		args...,
	); err != nil {
		return nil, err
	}
	out := make([]model.Coupon, 0, len(list))
	for i := range list {
		ds := couponDisplayStatus(&list[i], now)
		if ds != "active" && ds != "sold_out" {
			continue
		}
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
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+userCouponColumns+" FROM user_coupons WHERE user_id=? AND status=? ORDER BY valid_end ASC",
		userID, model.UserCouponUnused,
	)
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
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+couponColumns+" FROM coupons WHERE status=? AND ((issuer_type=? AND shop_id=0) OR (issuer_type=? AND shop_id=?))",
		model.CouponStatusOn, model.CouponIssuerPlatform, model.CouponIssuerShop, shopID,
	)
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
