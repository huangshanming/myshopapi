package repository

import (
	"context"
	"errors"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	slotPackageColumns = "id, slot_type, name, price, duration_days, status, sort, remark, created_at, updated_at"
	slotSettingColumns = "slot_type, home_limit, updated_at"
	slotOrderColumns = "id, shop_id, slot_type, package_id, target_id, amount, duration_days, start_at, end_at, status, pay_source, wallet_log_id, operator_id, created_at, updated_at"
)

func (r *MerchantRepository) ExpireDueSlotOrders(ctx context.Context) {
	now := time.Now()
	_, _ = r.conn.ExecCtx(ctx,
		"UPDATE homepage_slot_orders SET status=? WHERE status=? AND end_at<?",
		model.SlotOrderExpired, model.SlotOrderActive, now,
	)
}

func (r *MerchantRepository) ListSlotPackages(ctx context.Context, slotType string, onlyOn bool) ([]model.HomepageSlotPackage, error) {
	where := "1=1"
	args := make([]any, 0, 2)
	if slotType != "" {
		where += " AND slot_type=?"
		args = append(args, slotType)
	}
	if onlyOn {
		where += " AND status=?"
		args = append(args, model.SlotPkgOn)
	}
	var list []model.HomepageSlotPackage
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+slotPackageColumns+" FROM homepage_slot_packages WHERE "+where+" ORDER BY sort ASC, id ASC",
		args...,
	)
	return list, err
}

func (r *MerchantRepository) GetSlotPackage(ctx context.Context, id uint64) (*model.HomepageSlotPackage, error) {
	var p model.HomepageSlotPackage
	err := r.conn.QueryRowCtx(ctx, &p,
		"SELECT "+slotPackageColumns+" FROM homepage_slot_packages WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MerchantRepository) CreateSlotPackage(ctx context.Context, p *model.HomepageSlotPackage) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO homepage_slot_packages (slot_type, name, price, duration_days, status, sort, remark)
		 VALUES (?,?,?,?,?,?,?)`,
		p.SlotType, p.Name, p.Price, p.DurationDays, p.Status, p.Sort, p.Remark,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *MerchantRepository) UpdateSlotPackage(ctx context.Context, p *model.HomepageSlotPackage) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE homepage_slot_packages SET slot_type=?, name=?, price=?, duration_days=?, status=?, sort=?, remark=? WHERE id=?`,
		p.SlotType, p.Name, p.Price, p.DurationDays, p.Status, p.Sort, p.Remark, p.ID,
	)
	return err
}

func (r *MerchantRepository) ListSlotSettings(ctx context.Context) ([]model.HomepageSlotSetting, error) {
	var list []model.HomepageSlotSetting
	err := r.conn.QueryRowsCtx(ctx, &list, "SELECT "+slotSettingColumns+" FROM homepage_slot_settings")
	return list, err
}

func (r *MerchantRepository) GetSlotSetting(ctx context.Context, slotType string) (*model.HomepageSlotSetting, error) {
	var s model.HomepageSlotSetting
	err := r.conn.QueryRowCtx(ctx, &s,
		"SELECT "+slotSettingColumns+" FROM homepage_slot_settings WHERE slot_type=? LIMIT 1", slotType,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		s = model.HomepageSlotSetting{SlotType: slotType, HomeLimit: 8}
		_, _ = r.conn.ExecCtx(ctx,
			"INSERT INTO homepage_slot_settings (slot_type, home_limit) VALUES (?,?)", slotType, s.HomeLimit,
		)
		return &s, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) UpsertSlotSetting(ctx context.Context, slotType string, homeLimit int) error {
	var s model.HomepageSlotSetting
	err := r.conn.QueryRowCtx(ctx, &s,
		"SELECT "+slotSettingColumns+" FROM homepage_slot_settings WHERE slot_type=? LIMIT 1", slotType,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err := r.conn.ExecCtx(ctx,
			"INSERT INTO homepage_slot_settings (slot_type, home_limit) VALUES (?,?)", slotType, homeLimit,
		)
		return err
	}
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx,
		"UPDATE homepage_slot_settings SET home_limit=? WHERE slot_type=?", homeLimit, slotType,
	)
	return err
}

func (r *MerchantRepository) ListSlotOrders(ctx context.Context, shopID uint64, slotType, status string, page, pageSize int) ([]model.HomepageSlotOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	where := "1=1"
	args := make([]any, 0, 3)
	if shopID > 0 {
		where += " AND shop_id=?"
		args = append(args, shopID)
	}
	if slotType != "" {
		where += " AND slot_type=?"
		args = append(args, slotType)
	}
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM homepage_slot_orders WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.HomepageSlotOrder
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+slotOrderColumns+" FROM homepage_slot_orders WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) LatestActiveSlotOrder(ctx context.Context, shopID uint64, slotType string, targetID uint64) (*model.HomepageSlotOrder, error) {
	r.ExpireDueSlotOrders(ctx)
	args := []any{shopID, slotType, model.SlotOrderActive}
	query := "SELECT " + slotOrderColumns + " FROM homepage_slot_orders WHERE shop_id=? AND slot_type=? AND status=?"
	if slotType == model.SlotArticle {
		query += " AND target_id=?"
		args = append(args, targetID)
	} else {
		query += " AND target_id=?"
		args = append(args, shopID)
	}
	query += " ORDER BY end_at DESC LIMIT 1"
	var o model.HomepageSlotOrder
	err := r.conn.QueryRowCtx(ctx, &o, query, args...)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *MerchantRepository) PurchaseSlotOrder(ctx context.Context, order *model.HomepageSlotOrder, deductWallet bool, operatorID *uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		now := time.Now()
		_, _ = session.ExecCtx(ctx,
			"UPDATE homepage_slot_orders SET status=? WHERE status=? AND end_at<?",
			model.SlotOrderExpired, model.SlotOrderActive, now,
		)

		targetID := order.TargetID
		args := []any{order.ShopID, order.SlotType, model.SlotOrderActive}
		query := "SELECT " + slotOrderColumns + " FROM homepage_slot_orders WHERE shop_id=? AND slot_type=? AND status=?"
		if order.SlotType == model.SlotArticle {
			query += " AND target_id=?"
			args = append(args, targetID)
		} else {
			query += " AND target_id=?"
			args = append(args, order.ShopID)
		}
		query += " ORDER BY end_at DESC LIMIT 1"

		start := now
		var prev model.HomepageSlotOrder
		if err := session.QueryRowCtx(ctx, &prev, query, args...); err == nil {
			if time.Time(prev.EndAt).After(now) {
				start = time.Time(prev.EndAt)
			}
		} else if !errors.Is(err, sqlx.ErrNotFound) {
			return err
		}

		end := start.Add(time.Duration(order.DurationDays) * 24 * time.Hour)
		order.StartAt = common.LocalTime(start)
		order.EndAt = common.LocalTime(end)
		order.Status = model.SlotOrderActive
		if order.SlotType != model.SlotArticle {
			order.TargetID = order.ShopID
		}

		if deductWallet {
			w, err := r.lockShopWallet(ctx, session, order.ShopID)
			if err != nil {
				return err
			}
			if w.Balance+0.0001 < order.Amount {
				return errors.New("余额不足，请联系平台充值后再购买")
			}
			w.Balance -= order.Amount
			if _, err := session.ExecCtx(ctx,
				"UPDATE shop_wallets SET balance=? WHERE shop_id=?", w.Balance, order.ShopID,
			); err != nil {
				return err
			}
			res, err := session.ExecCtx(ctx,
				`INSERT INTO homepage_slot_orders (shop_id, slot_type, package_id, target_id, amount, duration_days, start_at, end_at, status, pay_source, operator_id)
				 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
				order.ShopID, order.SlotType, order.PackageID, order.TargetID, order.Amount, order.DurationDays,
				order.StartAt, order.EndAt, order.Status, order.PaySource, order.OperatorID,
			)
			if err != nil {
				return err
			}
			orderID, err := lastInsertID(res)
			if err != nil {
				return err
			}
			order.ID = orderID
			logRes, err := session.ExecCtx(ctx,
				`INSERT INTO shop_wallet_logs (shop_id, change_type, amount, balance_after, frozen_after, deposit_after, remark, operator_user_id, ref_type, ref_id)
				 VALUES (?,?,?,?,?,?,?,?,?,?)`,
				order.ShopID, model.WalletLogHomepageSlot, -order.Amount, w.Balance, w.FrozenBalance, w.Deposit,
				"首页展位套餐购买", operatorID, "homepage_slot_order", orderID,
			)
			if err != nil {
				return err
			}
			walletLogID, err := lastInsertID(logRes)
			if err != nil {
				return err
			}
			order.WalletLogID = walletLogID
			_, err = session.ExecCtx(ctx,
				"UPDATE homepage_slot_orders SET wallet_log_id=? WHERE id=?", walletLogID, orderID,
			)
			return err
		}
		res, err := session.ExecCtx(ctx,
			`INSERT INTO homepage_slot_orders (shop_id, slot_type, package_id, target_id, amount, duration_days, start_at, end_at, status, pay_source, operator_id)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			order.ShopID, order.SlotType, order.PackageID, order.TargetID, order.Amount, order.DurationDays,
			order.StartAt, order.EndAt, order.Status, order.PaySource, order.OperatorID,
		)
		if err != nil {
			return err
		}
		orderID, err := lastInsertID(res)
		if err != nil {
			return err
		}
		order.ID = orderID
		return nil
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
	if err := r.conn.QueryRowCtx(ctx, &total, countSQL,
		slotType, model.SlotOrderActive, now, now, model.ShopApproved,
	); err != nil {
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
SELECT s.id, s.name, s.logo, s.contact_name, s.contact_phone, s.description, s.category, s.province, s.city, s.district, s.address, s.business_license_no, s.legal_person, s.license_image, s.storefront_image, s.owner_user_id, s.status, s.reject_reason, s.created_at, s.updated_at
` + baseFrom + `
ORDER BY CASE WHEN o.target_id IS NULL THEN 0 ELSE 1 END DESC, s.id DESC
LIMIT ? OFFSET ?`

	var list []model.Shop
	if err := r.conn.QueryRowsCtx(ctx, &list, listSQL,
		slotType, model.SlotOrderActive, now, now, model.ShopApproved, limit, offset,
	); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *MerchantRepository) GetArticleTitle(ctx context.Context, id uint64) (string, error) {
	var title string
	err := r.conn.QueryRowCtx(ctx, &title,
		"SELECT title FROM community_article WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return "", err
	}
	if title == "" {
		return "", sqlx.ErrNotFound
	}
	return title, nil
}

func (r *MerchantRepository) ArticlePublishedForShop(ctx context.Context, articleID, shopID uint64) bool {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM community_article WHERE id=? AND shop_id=? AND status=? AND audit_status=?",
		articleID, shopID, "published", "approved",
	)
	return err == nil && n > 0
}

func (r *MerchantRepository) ActivePaidTargetIDs(ctx context.Context, slotType string) (map[uint64]bool, error) {
	r.ExpireDueSlotOrders(ctx)
	now := time.Now()
	var ids []uint64
	err := r.conn.QueryRowsCtx(ctx, &ids,
		"SELECT DISTINCT target_id FROM homepage_slot_orders WHERE slot_type=? AND status=? AND start_at<=? AND end_at>?",
		slotType, model.SlotOrderActive, now, now,
	)
	m := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m, err
}
