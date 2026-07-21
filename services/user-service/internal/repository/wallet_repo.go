package repository

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *UserRepository) EnsureWallet(ctx context.Context, userID uint64) (*model.UserWallet, error) {
	var w model.UserWallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w = model.UserWallet{UserID: userID}
		if err := r.db.WithContext(ctx).Create(&w).Error; err != nil {
			return nil, err
		}
		return &w, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *UserRepository) GetWallet(ctx context.Context, userID uint64) (*model.UserWallet, error) {
	return r.EnsureWallet(ctx, userID)
}

func (r *UserRepository) lockWallet(tx *gorm.DB, userID uint64) (*model.UserWallet, error) {
	var w model.UserWallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w = model.UserWallet{UserID: userID}
		if err := tx.Create(&w).Error; err != nil {
			return nil, err
		}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&w).Error
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *UserRepository) writeLog(tx *gorm.DB, w *model.UserWallet, changeType string, amount float64, remark string, op *uint64, refType string, refID uint64) error {
	return tx.Create(&model.UserWalletLog{
		UserID:         w.UserID,
		ChangeType:     changeType,
		Amount:         amount,
		BalanceAfter:   w.Balance,
		FrozenAfter:    w.FrozenBalance,
		Remark:         remark,
		OperatorUserID: op,
		RefType:        refType,
		RefID:          refID,
	}).Error
}

func (r *UserRepository) HasWalletLog(ctx context.Context, userID uint64, changeType, refType string, refID uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.UserWalletLog{}).
		Where("user_id = ? AND change_type = ? AND ref_type = ? AND ref_id = ?", userID, changeType, refType, refID).
		Count(&n).Error
	return n > 0, err
}

func (r *UserRepository) AdjustWallet(ctx context.Context, userID uint64, field string, amount float64, remark string, operatorID *uint64) (*model.UserWallet, error) {
	if field == "" {
		field = model.UserWalletFieldBalance
	}
	var out *model.UserWallet
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		w, err := r.lockWallet(tx, userID)
		if err != nil {
			return err
		}
		updates := map[string]interface{}{}
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
			updates["balance"] = w.Balance
		case model.UserWalletFieldFrozen:
			next := w.FrozenBalance + amount
			if next < -0.0001 {
				return errors.New("冻结余额不足")
			}
			if next < 0 {
				next = 0
			}
			w.FrozenBalance = next
			updates["frozen_balance"] = w.FrozenBalance
		default:
			return errors.New("不支持的调账字段")
		}
		if err := tx.Model(w).Updates(updates).Error; err != nil {
			return err
		}
		if err := r.writeLog(tx, w, model.UserWalletLogAdminAdjust, amount, remark, operatorID, field, 0); err != nil {
			return err
		}
		out = w
		return nil
	})
	return out, err
}

func (r *UserRepository) ListWalletLogs(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserWalletLog, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.UserWalletLog{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.UserWalletLog
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		w, err := r.lockWallet(tx, userID)
		if err != nil {
			return err
		}
		if w.Balance+0.0001 < amount {
			return errors.New("余额不足")
		}
		w.Balance -= amount
		w.FrozenBalance += amount
		if err := tx.Model(w).Updates(map[string]interface{}{
			"balance":        w.Balance,
			"frozen_balance": w.FrozenBalance,
		}).Error; err != nil {
			return err
		}
		remark := "下单冻结 " + orderNo
		return r.writeLog(tx, w, model.UserWalletLogOrderFreeze, -amount, remark, nil, "order", orderID)
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		w, err := r.lockWallet(tx, userID)
		if err != nil {
			return err
		}
		if w.FrozenBalance+0.0001 < amount {
			amount = w.FrozenBalance
		}
		w.FrozenBalance -= amount
		w.Balance += amount
		if err := tx.Model(w).Updates(map[string]interface{}{
			"balance":        w.Balance,
			"frozen_balance": w.FrozenBalance,
		}).Error; err != nil {
			return err
		}
		remark := "取消订单解冻 " + orderNo
		return r.writeLog(tx, w, model.UserWalletLogOrderUnfreeze, amount, remark, nil, "order", orderID)
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		w, err := r.lockWallet(tx, userID)
		if err != nil {
			return err
		}
		if w.FrozenBalance+0.0001 < amount {
			amount = w.FrozenBalance
		}
		w.FrozenBalance -= amount
		if err := tx.Model(w).Update("frozen_balance", w.FrozenBalance).Error; err != nil {
			return err
		}
		remark := "订单确认实扣 " + orderNo
		return r.writeLog(tx, w, model.UserWalletLogOrderSettle, -amount, remark, nil, "order", orderID)
	})
}
