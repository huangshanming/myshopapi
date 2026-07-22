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
	themeSlotColumns = "id, slot_key, position, name, `desc`, cover_url, default_link_type, default_link_id, status, sort, created_at, updated_at"
	themePackageColumns = "id, theme_slot_id, name, price, duration_days, status, sort, remark, created_at, updated_at"
	themeOrderColumns = "id, shop_id, theme_slot_id, package_id, title, subtitle, cover_url, link_type, link_id, amount, duration_days, start_at, end_at, status, pay_source, wallet_log_id, operator_id, created_at, updated_at"
)

func (r *MerchantRepository) ExpireDueThemeOrders(ctx context.Context) {
	now := time.Now()
	_, _ = r.conn.ExecCtx(ctx,
		"UPDATE homepage_theme_orders SET status=? WHERE status=? AND end_at<?",
		model.ThemeOrderExpired, model.ThemeOrderActive, now,
	)
}

func (r *MerchantRepository) ListThemeSlots(ctx context.Context, onlyOn bool) ([]model.HomepageThemeSlot, error) {
	r.ExpireDueThemeOrders(ctx)
	where := "1=1"
	args := make([]any, 0, 1)
	if onlyOn {
		where += " AND status=?"
		args = append(args, model.ThemeSlotOn)
	}
	var list []model.HomepageThemeSlot
	if err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+themeSlotColumns+" FROM homepage_theme_slots WHERE "+where+" ORDER BY position ASC, sort ASC, id ASC",
		args...,
	); err != nil {
		return nil, err
	}
	now := time.Now()
	for i := range list {
		var o model.HomepageThemeOrder
		e := r.conn.QueryRowCtx(ctx, &o,
			"SELECT "+themeOrderColumns+" FROM homepage_theme_orders WHERE theme_slot_id=? AND status=? AND end_at>? ORDER BY end_at DESC LIMIT 1",
			list[i].ID, model.ThemeOrderActive, now,
		)
		if e == nil {
			list[i].HasActive = true
			list[i].OccupiedUntil = time.Time(o.EndAt).Format("2006-01-02 15:04:05")
		}
	}
	return list, nil
}

func (r *MerchantRepository) GetThemeSlot(ctx context.Context, id uint64) (*model.HomepageThemeSlot, error) {
	var s model.HomepageThemeSlot
	err := r.conn.QueryRowCtx(ctx, &s,
		"SELECT "+themeSlotColumns+" FROM homepage_theme_slots WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) UpdateThemeSlot(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("homepage_theme_slots", updates, "id=?", id)
	if err != nil {
		return err
	}
	n, err := execRows(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("坑位不存在")
	}
	return nil
}

func (r *MerchantRepository) ListThemePackages(ctx context.Context, themeSlotID uint64, onlyOn bool) ([]model.HomepageThemePackage, error) {
	where := "1=1"
	args := make([]any, 0, 2)
	if onlyOn {
		where += " AND status=?"
		args = append(args, model.ThemeSlotOn)
	}
	if themeSlotID > 0 {
		where += " AND (theme_slot_id=0 OR theme_slot_id=?)"
		args = append(args, themeSlotID)
	}
	var list []model.HomepageThemePackage
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+themePackageColumns+" FROM homepage_theme_packages WHERE "+where+" ORDER BY sort ASC, id ASC",
		args...,
	)
	return list, err
}

func (r *MerchantRepository) GetThemePackage(ctx context.Context, id uint64) (*model.HomepageThemePackage, error) {
	var p model.HomepageThemePackage
	err := r.conn.QueryRowCtx(ctx, &p,
		"SELECT "+themePackageColumns+" FROM homepage_theme_packages WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MerchantRepository) CreateThemePackage(ctx context.Context, p *model.HomepageThemePackage) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO homepage_theme_packages (theme_slot_id, name, price, duration_days, status, sort, remark)
		 VALUES (?,?,?,?,?,?,?)`,
		p.ThemeSlotID, p.Name, p.Price, p.DurationDays, p.Status, p.Sort, p.Remark,
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

func (r *MerchantRepository) UpdateThemePackage(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("homepage_theme_packages", updates, "id=?", id)
	if err != nil {
		return err
	}
	n, err := execRows(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
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
	where := "1=1"
	args := make([]any, 0, 2)
	if shopID > 0 {
		where += " AND shop_id=?"
		args = append(args, shopID)
	}
	if themeSlotID > 0 {
		where += " AND theme_slot_id=?"
		args = append(args, themeSlotID)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM homepage_theme_orders WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.HomepageThemeOrder
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+themeOrderColumns+" FROM homepage_theme_orders WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) ActiveThemeOrderForSlot(ctx context.Context, themeSlotID uint64) (*model.HomepageThemeOrder, error) {
	r.ExpireDueThemeOrders(ctx)
	now := time.Now()
	var o model.HomepageThemeOrder
	err := r.conn.QueryRowCtx(ctx, &o,
		"SELECT "+themeOrderColumns+" FROM homepage_theme_orders WHERE theme_slot_id=? AND status=? AND start_at<=? AND end_at>? ORDER BY end_at DESC LIMIT 1",
		themeSlotID, model.ThemeOrderActive, now, now,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
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
	err := r.conn.QueryRowCtx(ctx, &o,
		"SELECT "+themeOrderColumns+" FROM homepage_theme_orders WHERE theme_slot_id=? AND status=? ORDER BY end_at DESC LIMIT 1",
		themeSlotID, model.ThemeOrderActive,
	)
	if err != nil {
		return now
	}
	end := time.Time(o.EndAt)
	if end.After(now) {
		return end
	}
	return now
}

func (r *MerchantRepository) PurchaseThemeOrder(ctx context.Context, order *model.HomepageThemeOrder, deductWallet bool, operatorID *uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		now := time.Now()
		_, _ = session.ExecCtx(ctx,
			"UPDATE homepage_theme_orders SET status=? WHERE status=? AND end_at<?",
			model.ThemeOrderExpired, model.ThemeOrderActive, now,
		)

		start := now
		var prev model.HomepageThemeOrder
		if err := session.QueryRowCtx(ctx, &prev,
			"SELECT "+themeOrderColumns+" FROM homepage_theme_orders WHERE theme_slot_id=? AND status=? ORDER BY end_at DESC LIMIT 1",
			order.ThemeSlotID, model.ThemeOrderActive,
		); err == nil {
			if time.Time(prev.EndAt).After(now) {
				start = time.Time(prev.EndAt)
			}
		} else if !errors.Is(err, sqlx.ErrNotFound) {
			return err
		}

		end := start.Add(time.Duration(order.DurationDays) * 24 * time.Hour)
		order.StartAt = common.LocalTime(start)
		order.EndAt = common.LocalTime(end)
		order.Status = model.ThemeOrderActive

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
				`INSERT INTO homepage_theme_orders (shop_id, theme_slot_id, package_id, title, subtitle, cover_url, link_type, link_id, amount, duration_days, start_at, end_at, status, pay_source, operator_id)
				 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
				order.ShopID, order.ThemeSlotID, order.PackageID, order.Title, order.Subtitle, order.CoverURL,
				order.LinkType, order.LinkID, order.Amount, order.DurationDays, order.StartAt, order.EndAt,
				order.Status, order.PaySource, order.OperatorID,
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
				order.ShopID, model.WalletLogThemeSlot, -order.Amount, w.Balance, w.FrozenBalance, w.Deposit,
				"主题集市坑位购买", operatorID, "homepage_theme_order", orderID,
			)
			if err != nil {
				return err
			}
			logID, err := lastInsertID(logRes)
			if err != nil {
				return err
			}
			order.WalletLogID = logID
			_, err = session.ExecCtx(ctx,
				"UPDATE homepage_theme_orders SET wallet_log_id=? WHERE id=?", logID, orderID,
			)
			return err
		}
		res, err := session.ExecCtx(ctx,
			`INSERT INTO homepage_theme_orders (shop_id, theme_slot_id, package_id, title, subtitle, cover_url, link_type, link_id, amount, duration_days, start_at, end_at, status, pay_source, operator_id)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			order.ShopID, order.ThemeSlotID, order.PackageID, order.Title, order.Subtitle, order.CoverURL,
			order.LinkType, order.LinkID, order.Amount, order.DurationDays, order.StartAt, order.EndAt,
			order.Status, order.PaySource, order.OperatorID,
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

func (r *MerchantRepository) CategoryExistsShow(ctx context.Context, id uint64) bool {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM product_categories WHERE id=? AND is_show=1", id,
	)
	return err == nil && n > 0
}

func (r *MerchantRepository) ProductOnSaleOfShop(ctx context.Context, productID, shopID uint64) bool {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM products WHERE id=? AND shop_id=? AND status=?", productID, shopID, "on_sale",
	)
	return err == nil && n > 0
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
		e := r.conn.QueryRowCtx(ctx, &o,
			"SELECT "+themeOrderColumns+" FROM homepage_theme_orders WHERE theme_slot_id=? AND status=? AND start_at<=? AND end_at>? ORDER BY end_at DESC LIMIT 1",
			s.ID, model.ThemeOrderActive, now, now,
		)
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
