package repository

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const walletColumns = "user_id, balance, frozen_balance, created_at, updated_at"

func (r *UserRepository) EnsureWallet(ctx context.Context, userID uint64) (*model.UserWallet, error) {
	var w model.UserWallet
	err := r.conn.QueryRowPartialCtx(ctx, &w,
		"SELECT "+walletColumns+" FROM user_wallets WHERE user_id=? LIMIT 1", userID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = r.conn.ExecCtx(ctx,
			"INSERT INTO user_wallets (user_id, balance, frozen_balance) VALUES (?, 0, 0)", userID,
		)
		if err != nil {
			return nil, err
		}
		err = r.conn.QueryRowPartialCtx(ctx, &w,
			"SELECT "+walletColumns+" FROM user_wallets WHERE user_id=? LIMIT 1", userID,
		)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *UserRepository) GetWallet(ctx context.Context, userID uint64) (*model.UserWallet, error) {
	return r.EnsureWallet(ctx, userID)
}

func (r *UserRepository) lockWallet(ctx context.Context, session sqlx.Session, userID uint64) (*model.UserWallet, error) {
	var w model.UserWallet
	err := session.QueryRowPartialCtx(ctx, &w,
		"SELECT "+walletColumns+" FROM user_wallets WHERE user_id=? FOR UPDATE", userID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = session.ExecCtx(ctx,
			"INSERT INTO user_wallets (user_id, balance, frozen_balance) VALUES (?, 0, 0)", userID,
		)
		if err != nil {
			return nil, err
		}
		err = session.QueryRowPartialCtx(ctx, &w,
			"SELECT "+walletColumns+" FROM user_wallets WHERE user_id=? FOR UPDATE", userID,
		)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *UserRepository) writeLog(ctx context.Context, session sqlx.Session, w *model.UserWallet, changeType string, amount float64, remark string, op *uint64, refType string, refID uint64) error {
	_, err := session.ExecCtx(ctx,
		`INSERT INTO user_wallet_logs (user_id, change_type, amount, balance_after, frozen_after, remark, operator_user_id, ref_type, ref_id)
		 VALUES (?,?,?,?,?,?,?,?,?)`,
		w.UserID, changeType, amount, w.Balance, w.FrozenBalance, remark, op, refType, refID,
	)
	return err
}

func (r *UserRepository) HasWalletLog(ctx context.Context, userID uint64, changeType, refType string, refID uint64) (bool, error) {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_wallet_logs WHERE user_id=? AND change_type=? AND ref_type=? AND ref_id=?",
		userID, changeType, refType, refID,
	)
	return n > 0, err
}

func (r *UserRepository) AdjustWallet(ctx context.Context, userID uint64, field string, amount float64, remark string, operatorID *uint64) (*model.UserWallet, error) {
	if field == "" {
		field = model.UserWalletFieldBalance
	}
	var out *model.UserWallet
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		w, err := r.lockWallet(ctx, session, userID)
		if err != nil {
			return err
		}
		switch field {
		case model.UserWalletFieldBalance:
			next := w.Balance + amount
			if next < -0.0001 {
				return errors.New("可用余额不足")
			}
			if next < 0 {
				next = 0
			}
			w.Balance = next
		case model.UserWalletFieldFrozen:
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
			"UPDATE user_wallets SET balance=?, frozen_balance=? WHERE user_id=?",
			w.Balance, w.FrozenBalance, w.UserID,
		); err != nil {
			return err
		}
		if err := r.writeLog(ctx, session, w, model.UserWalletLogAdminAdjust, amount, remark, operatorID, field, 0); err != nil {
			return err
		}
		out = w
		return nil
	})
	return out, err
}

func (r *UserRepository) ListWalletLogs(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserWalletLog, int64, error) {
	total, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM user_wallet_logs WHERE user_id=?", userID,
	)
	if err != nil {
		return nil, 0, err
	}
	var list []model.UserWalletLog
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT id, user_id, change_type, amount, balance_after, frozen_after, IFNULL(remark,'') AS remark, IFNULL(operator_user_id,0) AS operator_user_id, IFNULL(ref_type,'') AS ref_type, ref_id, created_at "+
			"FROM user_wallet_logs WHERE user_id=? ORDER BY id DESC LIMIT ? OFFSET ?",
		userID, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *UserRepository) FreezeForOrder(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	if amount <= 0 {
		return errors.New("冻结金额无效")
	}
	exists, err := r.HasWalletLog(ctx, userID, model.UserWalletLogOrderFreeze, "order", orderID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		w, err := r.lockWallet(ctx, session, userID)
		if err != nil {
			return err
		}
		if w.Balance+0.0001 < amount {
			return errors.New("余额不足")
		}
		w.Balance -= amount
		w.FrozenBalance += amount
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_wallets SET balance=?, frozen_balance=? WHERE user_id=?",
			w.Balance, w.FrozenBalance, w.UserID,
		); err != nil {
			return err
		}
		remark := "下单冻结 " + orderNo
		return r.writeLog(ctx, session, w, model.UserWalletLogOrderFreeze, -amount, remark, nil, "order", orderID)
	})
}

func (r *UserRepository) UnfreezeOrder(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	if amount <= 0 {
		return nil
	}
	settled, err := r.HasWalletLog(ctx, userID, model.UserWalletLogOrderSettle, "order", orderID)
	if err != nil {
		return err
	}
	if settled {
		return nil
	}
	exists, err := r.HasWalletLog(ctx, userID, model.UserWalletLogOrderUnfreeze, "order", orderID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	frozen, err := r.HasWalletLog(ctx, userID, model.UserWalletLogOrderFreeze, "order", orderID)
	if err != nil {
		return err
	}
	if !frozen {
		return nil
	}
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		w, err := r.lockWallet(ctx, session, userID)
		if err != nil {
			return err
		}
		if w.FrozenBalance+0.0001 < amount {
			amount = w.FrozenBalance
		}
		w.FrozenBalance -= amount
		w.Balance += amount
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_wallets SET balance=?, frozen_balance=? WHERE user_id=?",
			w.Balance, w.FrozenBalance, w.UserID,
		); err != nil {
			return err
		}
		remark := "取消订单解冻 " + orderNo
		return r.writeLog(ctx, session, w, model.UserWalletLogOrderUnfreeze, amount, remark, nil, "order", orderID)
	})
}

func (r *UserRepository) SettleOrder(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	if amount <= 0 {
		return nil
	}
	exists, err := r.HasWalletLog(ctx, userID, model.UserWalletLogOrderSettle, "order", orderID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	unfrozen, err := r.HasWalletLog(ctx, userID, model.UserWalletLogOrderUnfreeze, "order", orderID)
	if err != nil {
		return err
	}
	if unfrozen {
		return nil
	}
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		w, err := r.lockWallet(ctx, session, userID)
		if err != nil {
			return err
		}
		if w.FrozenBalance+0.0001 < amount {
			amount = w.FrozenBalance
		}
		w.FrozenBalance -= amount
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_wallets SET frozen_balance=? WHERE user_id=?",
			w.FrozenBalance, w.UserID,
		); err != nil {
			return err
		}
		remark := "订单确认实扣 " + orderNo
		return r.writeLog(ctx, session, w, model.UserWalletLogOrderSettle, -amount, remark, nil, "order", orderID)
	})
}
