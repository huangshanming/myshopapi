package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	shopWalletColumns = "shop_id, balance, frozen_balance, deposit, created_at, updated_at"
	seckillRuleColumns = "id, duration_hours, apply_fee, max_entries_per_shop, status, created_at, updated_at"
	seckillSessionColumns = "id, rule_id, start_at, end_at, status, created_at, updated_at"
	seckillEntryColumns = "id, session_id, shop_id, product_id, product_name, product_image, origin_price, seckill_price, seckill_stock, fee_amount, status, auto_renew, created_at, updated_at"
)

func (r *MerchantRepository) lockShopWallet(ctx context.Context, session sqlx.Session, shopID uint64) (*model.ShopWallet, error) {
	var w model.ShopWallet
	err := session.QueryRowCtx(ctx, &w,
		"SELECT "+shopWalletColumns+" FROM shop_wallets WHERE shop_id=? FOR UPDATE", shopID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		if _, err = session.ExecCtx(ctx,
			"INSERT INTO shop_wallets (shop_id, balance, frozen_balance, deposit) VALUES (?, 0, 0, 0)", shopID,
		); err != nil {
			return nil, err
		}
		err = session.QueryRowCtx(ctx, &w,
			"SELECT "+shopWalletColumns+" FROM shop_wallets WHERE shop_id=? FOR UPDATE", shopID,
		)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *MerchantRepository) writeShopWalletLog(ctx context.Context, session sqlx.Session, w *model.ShopWallet, changeType string, amount float64, remark string, operatorID *uint64, refType string, refID uint64) error {
	_, err := session.ExecCtx(ctx,
		`INSERT INTO shop_wallet_logs (shop_id, change_type, amount, balance_after, frozen_after, deposit_after, remark, operator_user_id, ref_type, ref_id)
		 VALUES (?,?,?,?,?,?,?,?,?,?)`,
		w.ShopID, changeType, amount, w.Balance, w.FrozenBalance, w.Deposit, remark, operatorID, refType, refID,
	)
	return err
}

func (r *MerchantRepository) EnsureWallet(ctx context.Context, shopID uint64) (*model.ShopWallet, error) {
	var w model.ShopWallet
	err := r.conn.QueryRowCtx(ctx, &w,
		"SELECT "+shopWalletColumns+" FROM shop_wallets WHERE shop_id=? LIMIT 1", shopID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = r.conn.ExecCtx(ctx,
			"INSERT INTO shop_wallets (shop_id, balance, frozen_balance, deposit) VALUES (?, 0, 0, 0)", shopID,
		)
		if err != nil {
			return nil, err
		}
		err = r.conn.QueryRowCtx(ctx, &w,
			"SELECT "+shopWalletColumns+" FROM shop_wallets WHERE shop_id=? LIMIT 1", shopID,
		)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *MerchantRepository) GetWallet(ctx context.Context, shopID uint64) (*model.ShopWallet, error) {
	return r.EnsureWallet(ctx, shopID)
}

func (r *MerchantRepository) AdjustWallet(ctx context.Context, shopID uint64, field string, amount float64, changeType, remark string, operatorID *uint64, refType string, refID uint64) (*model.ShopWallet, error) {
	if field == "" {
		field = model.WalletFieldBalance
	}
	var out *model.ShopWallet
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		w, err := r.lockShopWallet(ctx, session, shopID)
		if err != nil {
			return err
		}
		switch field {
		case model.WalletFieldBalance:
			next := w.Balance + amount
			if next < -0.0001 {
				return errors.New("可用余额不足")
			}
			if next < 0 {
				next = 0
			}
			w.Balance = next
		case model.WalletFieldDeposit:
			next := w.Deposit + amount
			if next < -0.0001 {
				return errors.New("保证金不足")
			}
			if next < 0 {
				next = 0
			}
			w.Deposit = next
		case model.WalletFieldFrozen:
			next := w.FrozenBalance + amount
			if next < -0.0001 {
				return errors.New("冻结余额不足")
			}
			if next < 0 {
				next = 0
			}
			w.FrozenBalance = next
		default:
			return errors.New("不支持的调账字段")
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE shop_wallets SET balance=?, frozen_balance=?, deposit=? WHERE shop_id=?",
			w.Balance, w.FrozenBalance, w.Deposit, shopID,
		); err != nil {
			return err
		}
		if err := r.writeShopWalletLog(ctx, session, w, changeType, amount, remark, operatorID, refType, refID); err != nil {
			return err
		}
		out = w
		return nil
	})
	return out, err
}

func (r *MerchantRepository) ListWalletLogs(ctx context.Context, shopID uint64, page, pageSize int) ([]model.ShopWalletLog, int64, error) {
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM shop_wallet_logs WHERE shop_id=?", shopID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.ShopWalletLog
	err = r.conn.QueryRowsCtx(ctx, &list,
		`SELECT id, shop_id, change_type, amount, balance_after, frozen_after, deposit_after, remark, operator_user_id, ref_type, ref_id, created_at
		 FROM shop_wallet_logs WHERE shop_id=? ORDER BY id DESC LIMIT ? OFFSET ?`,
		shopID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) GetActiveSeckillRule(ctx context.Context) (*model.SeckillRule, error) {
	var rule model.SeckillRule
	err := r.conn.QueryRowCtx(ctx, &rule,
		"SELECT "+seckillRuleColumns+" FROM seckill_rules WHERE status=? ORDER BY id DESC LIMIT 1",
		model.SeckillRuleOn,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		rule = model.SeckillRule{
			DurationHours:     24,
			ApplyFee:          10,
			MaxEntriesPerShop: 5,
			Status:            model.SeckillRuleOn,
		}
		res, err := r.conn.ExecCtx(ctx,
			"INSERT INTO seckill_rules (duration_hours, apply_fee, max_entries_per_shop, status) VALUES (?,?,?,?)",
			rule.DurationHours, rule.ApplyFee, rule.MaxEntriesPerShop, rule.Status,
		)
		if err != nil {
			return nil, err
		}
		id, err := lastInsertID(res)
		if err != nil {
			return nil, err
		}
		rule.ID = id
		return &rule, nil
	}
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *MerchantRepository) SaveSeckillRule(ctx context.Context, rule *model.SeckillRule) error {
	if rule.ID > 0 {
		_, err := r.conn.ExecCtx(ctx,
			"UPDATE seckill_rules SET duration_hours=?, apply_fee=?, max_entries_per_shop=?, status=? WHERE id=?",
			rule.DurationHours, rule.ApplyFee, rule.MaxEntriesPerShop, rule.Status, rule.ID,
		)
		return err
	}
	res, err := r.conn.ExecCtx(ctx,
		"INSERT INTO seckill_rules (duration_hours, apply_fee, max_entries_per_shop, status) VALUES (?,?,?,?)",
		rule.DurationHours, rule.ApplyFee, rule.MaxEntriesPerShop, rule.Status,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	rule.ID = id
	return nil
}

func (r *MerchantRepository) GetActiveSession(ctx context.Context) (*model.SeckillSession, error) {
	var s model.SeckillSession
	err := r.conn.QueryRowCtx(ctx, &s,
		"SELECT "+seckillSessionColumns+" FROM seckill_sessions WHERE status=? ORDER BY id DESC LIMIT 1",
		model.SeckillSessionActive,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) CreateSession(ctx context.Context, ruleID uint64, start, end time.Time) (*model.SeckillSession, error) {
	s := &model.SeckillSession{
		RuleID:  ruleID,
		StartAt: common.LocalTime(start),
		EndAt:   common.LocalTime(end),
		Status:  model.SeckillSessionActive,
	}
	res, err := r.conn.ExecCtx(ctx,
		"INSERT INTO seckill_sessions (rule_id, start_at, end_at, status) VALUES (?,?,?,?)",
		s.RuleID, s.StartAt, s.EndAt, s.Status,
	)
	if err != nil {
		return nil, err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return nil, err
	}
	s.ID = id
	return s, nil
}

func (r *MerchantRepository) EndSession(ctx context.Context, id uint64) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE seckill_sessions SET status=? WHERE id=?", model.SeckillSessionEnded, id,
	)
	return err
}

func (r *MerchantRepository) ListSessions(ctx context.Context, page, pageSize int) ([]model.SeckillSession, int64, error) {
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM seckill_sessions")
	if err != nil {
		return nil, 0, err
	}
	var list []model.SeckillSession
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+seckillSessionColumns+" FROM seckill_sessions ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) FindSession(ctx context.Context, id uint64) (*model.SeckillSession, error) {
	var s model.SeckillSession
	err := r.conn.QueryRowCtx(ctx, &s,
		"SELECT "+seckillSessionColumns+" FROM seckill_sessions WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) CountShopEntries(ctx context.Context, sessionID, shopID uint64) (int64, error) {
	return countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM seckill_entries WHERE session_id=? AND shop_id=? AND status=?",
		sessionID, shopID, model.SeckillEntryActive,
	)
}

func (r *MerchantRepository) CreateSeckillEntry(ctx context.Context, entry *model.SeckillEntry) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO seckill_entries (session_id, shop_id, product_id, product_name, product_image, origin_price, seckill_price, seckill_stock, fee_amount, status, auto_renew)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		entry.SessionID, entry.ShopID, entry.ProductID, entry.ProductName, entry.ProductImage,
		entry.OriginPrice, entry.SeckillPrice, entry.SeckillStock, entry.FeeAmount, entry.Status, entry.AutoRenew,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	entry.ID = id
	return nil
}

func (r *MerchantRepository) ApplySeckillEntry(ctx context.Context, entry *model.SeckillEntry, fee float64, operatorID *uint64) error {
	return r.applySeckillEntryTx(ctx, entry, fee, operatorID, model.WalletLogSeckillApply, "秒杀报名扣费")
}

func (r *MerchantRepository) applySeckillEntryTx(ctx context.Context, entry *model.SeckillEntry, fee float64, operatorID *uint64, changeType, remark string) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		w, err := r.lockShopWallet(ctx, session, entry.ShopID)
		if err != nil {
			return err
		}
		if w.Balance+0.0001 < fee {
			return errors.New("余额不足，请联系平台充值后再报名")
		}
		w.Balance -= fee
		if _, err := session.ExecCtx(ctx,
			"UPDATE shop_wallets SET balance=? WHERE shop_id=?", w.Balance, entry.ShopID,
		); err != nil {
			return err
		}
		entry.FeeAmount = fee
		entry.Status = model.SeckillEntryActive
		res, err := session.ExecCtx(ctx,
			`INSERT INTO seckill_entries (session_id, shop_id, product_id, product_name, product_image, origin_price, seckill_price, seckill_stock, fee_amount, status, auto_renew)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			entry.SessionID, entry.ShopID, entry.ProductID, entry.ProductName, entry.ProductImage,
			entry.OriginPrice, entry.SeckillPrice, entry.SeckillStock, entry.FeeAmount, entry.Status, entry.AutoRenew,
		)
		if err != nil {
			if isDup(err) {
				return errors.New("该商品已报名本场次")
			}
			return err
		}
		entryID, err := lastInsertID(res)
		if err != nil {
			return err
		}
		entry.ID = entryID
		return r.writeShopWalletLog(ctx, session, w, changeType, -fee, remark, operatorID, "seckill_entry", entryID)
	})
}

func (r *MerchantRepository) ListAutoRenewEntries(ctx context.Context, sessionID uint64) ([]model.SeckillEntry, error) {
	var list []model.SeckillEntry
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE session_id=? AND status=? AND auto_renew=1 ORDER BY id ASC",
		sessionID, model.SeckillEntryActive,
	)
	return list, err
}

func (r *MerchantRepository) SetSeckillAutoRenew(ctx context.Context, shopID, entryID uint64, autoRenew int8) error {
	if autoRenew != 0 {
		autoRenew = 1
	}
	n, err := execRows(ctx, r.conn,
		"UPDATE seckill_entries SET auto_renew=? WHERE id=? AND shop_id=?", autoRenew, entryID, shopID,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("报名记录不存在")
	}
	return nil
}

func (r *MerchantRepository) FindShopEntry(ctx context.Context, shopID, entryID uint64) (*model.SeckillEntry, error) {
	var e model.SeckillEntry
	err := r.conn.QueryRowCtx(ctx, &e,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE id=? AND shop_id=? LIMIT 1",
		entryID, shopID,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *MerchantRepository) RenewSeckillEntry(ctx context.Context, old *model.SeckillEntry, newSessionID uint64, fee float64, maxPerShop int) error {
	cnt, err := r.CountShopEntries(ctx, newSessionID, old.ShopID)
	if err != nil {
		return err
	}
	if maxPerShop > 0 && int(cnt) >= maxPerShop {
		return errors.New("已达本场次报名上限")
	}
	entry := &model.SeckillEntry{
		SessionID:    newSessionID,
		ShopID:       old.ShopID,
		ProductID:    old.ProductID,
		ProductName:  old.ProductName,
		ProductImage: old.ProductImage,
		OriginPrice:  old.OriginPrice,
		SeckillPrice: old.SeckillPrice,
		SeckillStock: old.SeckillStock,
		AutoRenew:    1,
	}
	return r.applySeckillEntryTx(ctx, entry, fee, nil, model.WalletLogSeckillRenew, "秒杀到期自动续费")
}

func isDup(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate") || strings.Contains(msg, "UNIQUE")
}

func (r *MerchantRepository) ListAdminEntries(ctx context.Context, sessionID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	where := "1=1"
	args := make([]any, 0, 1)
	if sessionID > 0 {
		where += " AND session_id=?"
		args = append(args, sessionID)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM seckill_entries WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.SeckillEntry
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) ListShopEntries(ctx context.Context, shopID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM seckill_entries WHERE shop_id=?", shopID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.SeckillEntry
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE shop_id=? ORDER BY id DESC LIMIT ? OFFSET ?",
		shopID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) ListActiveEntries(ctx context.Context, sessionID uint64, limit int) ([]model.SeckillEntry, error) {
	var list []model.SeckillEntry
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE session_id=? AND status=? ORDER BY id ASC LIMIT ?",
		sessionID, model.SeckillEntryActive, limit,
	)
	return list, err
}

func (r *MerchantRepository) ListActiveEntriesPage(ctx context.Context, sessionID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM seckill_entries WHERE session_id=? AND status=?",
		sessionID, model.SeckillEntryActive,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.SeckillEntry
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE session_id=? AND status=? ORDER BY id ASC LIMIT ? OFFSET ?",
		sessionID, model.SeckillEntryActive, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) FindSeckillEntry(ctx context.Context, id uint64) (*model.SeckillEntry, error) {
	var e model.SeckillEntry
	err := r.conn.QueryRowCtx(ctx, &e,
		"SELECT "+seckillEntryColumns+" FROM seckill_entries WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *MerchantRepository) DecrSeckillStock(ctx context.Context, entryID uint64, qty int) error {
	n, err := execRows(ctx, r.conn,
		`UPDATE seckill_entries SET seckill_stock=seckill_stock-? WHERE id=? AND status=? AND seckill_stock>=?`,
		qty, entryID, model.SeckillEntryActive, qty,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("秒杀库存不足或场次无效")
	}
	return nil
}

func (r *MerchantRepository) IncrSeckillStock(ctx context.Context, entryID uint64, qty int) error {
	if qty <= 0 {
		return nil
	}
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE seckill_entries SET seckill_stock=seckill_stock+? WHERE id=?", qty, entryID,
	)
	return err
}
