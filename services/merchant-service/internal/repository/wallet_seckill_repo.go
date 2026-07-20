package repository

import (
	"errors"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/merchant-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *MerchantRepository) EnsureWallet(shopID uint64) (*model.ShopWallet, error) {
	var w model.ShopWallet
	err := r.db.Where("shop_id = ?", shopID).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w = model.ShopWallet{ShopID: shopID}
		if err := r.db.Create(&w).Error; err != nil {
			return nil, err
		}
		return &w, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *MerchantRepository) GetWallet(shopID uint64) (*model.ShopWallet, error) {
	return r.EnsureWallet(shopID)
}

func (r *MerchantRepository) AdjustWallet(shopID uint64, field string, amount float64, changeType, remark string, operatorID *uint64, refType string, refID uint64) (*model.ShopWallet, error) {
	if field == "" {
		field = model.WalletFieldBalance
	}
	var out *model.ShopWallet
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var w model.ShopWallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("shop_id = ?", shopID).First(&w).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w = model.ShopWallet{ShopID: shopID}
			if err := tx.Create(&w).Error; err != nil {
				return err
			}
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("shop_id = ?", shopID).First(&w).Error
		}
		if err != nil {
			return err
		}
		updates := map[string]interface{}{}
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
			updates["balance"] = w.Balance
		case model.WalletFieldDeposit:
			next := w.Deposit + amount
			if next < -0.0001 {
				return errors.New("保证金不足")
			}
			if next < 0 {
				next = 0
			}
			w.Deposit = next
			updates["deposit"] = w.Deposit
		case model.WalletFieldFrozen:
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
		if err := tx.Model(&w).Updates(updates).Error; err != nil {
			return err
		}
		log := model.ShopWalletLog{
			ShopID:         shopID,
			ChangeType:     changeType,
			Amount:         amount,
			BalanceAfter:   w.Balance,
			FrozenAfter:    w.FrozenBalance,
			DepositAfter:   w.Deposit,
			Remark:         remark,
			OperatorUserID: operatorID,
			RefType:        refType,
			RefID:          refID,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		out = &w
		return nil
	})
	return out, err
}

func (r *MerchantRepository) ListWalletLogs(shopID uint64, page, pageSize int) ([]model.ShopWalletLog, int64, error) {
	q := r.db.Model(&model.ShopWalletLog{}).Where("shop_id = ?", shopID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.ShopWalletLog
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) GetActiveSeckillRule() (*model.SeckillRule, error) {
	var rule model.SeckillRule
	err := r.db.Where("status = ?", model.SeckillRuleOn).Order("id DESC").First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rule = model.SeckillRule{
			DurationHours:     24,
			ApplyFee:          10,
			MaxEntriesPerShop: 5,
			Status:            model.SeckillRuleOn,
		}
		if err := r.db.Create(&rule).Error; err != nil {
			return nil, err
		}
		return &rule, nil
	}
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *MerchantRepository) SaveSeckillRule(rule *model.SeckillRule) error {
	return r.db.Save(rule).Error
}

func (r *MerchantRepository) GetActiveSession() (*model.SeckillSession, error) {
	var s model.SeckillSession
	err := r.db.Where("status = ?", model.SeckillSessionActive).Order("id DESC").First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) CreateSession(ruleID uint64, start, end time.Time) (*model.SeckillSession, error) {
	s := &model.SeckillSession{
		RuleID:  ruleID,
		StartAt: common.LocalTime(start),
		EndAt:   common.LocalTime(end),
		Status:  model.SeckillSessionActive,
	}
	if err := r.db.Create(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *MerchantRepository) EndSession(id uint64) error {
	return r.db.Model(&model.SeckillSession{}).Where("id = ?", id).Update("status", model.SeckillSessionEnded).Error
}

func (r *MerchantRepository) ListSessions(page, pageSize int) ([]model.SeckillSession, int64, error) {
	q := r.db.Model(&model.SeckillSession{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SeckillSession
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) FindSession(id uint64) (*model.SeckillSession, error) {
	var s model.SeckillSession
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MerchantRepository) CountShopEntries(sessionID, shopID uint64) (int64, error) {
	var n int64
	err := r.db.Model(&model.SeckillEntry{}).
		Where("session_id = ? AND shop_id = ? AND status = ?", sessionID, shopID, model.SeckillEntryActive).
		Count(&n).Error
	return n, err
}

func (r *MerchantRepository) CreateSeckillEntry(entry *model.SeckillEntry) error {
	return r.db.Create(entry).Error
}

func (r *MerchantRepository) ApplySeckillEntry(entry *model.SeckillEntry, fee float64, operatorID *uint64) error {
	return r.applySeckillEntryTx(entry, fee, operatorID, model.WalletLogSeckillApply, "秒杀报名扣费")
}

func (r *MerchantRepository) applySeckillEntryTx(entry *model.SeckillEntry, fee float64, operatorID *uint64, changeType, remark string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var w model.ShopWallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("shop_id = ?", entry.ShopID).First(&w).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w = model.ShopWallet{ShopID: entry.ShopID}
			if err := tx.Create(&w).Error; err != nil {
				return err
			}
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("shop_id = ?", entry.ShopID).First(&w).Error
		}
		if err != nil {
			return err
		}
		if w.Balance+0.0001 < fee {
			return errors.New("余额不足，请联系平台充值后再报名")
		}
		w.Balance -= fee
		if err := tx.Model(&w).Update("balance", w.Balance).Error; err != nil {
			return err
		}
		entry.FeeAmount = fee
		entry.Status = model.SeckillEntryActive
		if err := tx.Create(entry).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) || isDup(err) {
				return errors.New("该商品已报名本场次")
			}
			return err
		}
		log := model.ShopWalletLog{
			ShopID:         entry.ShopID,
			ChangeType:     changeType,
			Amount:         -fee,
			BalanceAfter:   w.Balance,
			FrozenAfter:    w.FrozenBalance,
			DepositAfter:   w.Deposit,
			Remark:         remark,
			OperatorUserID: operatorID,
			RefType:        "seckill_entry",
			RefID:          entry.ID,
		}
		return tx.Create(&log).Error
	})
}

// ListAutoRenewEntries 到期场次中开启自动续费的报名（按报名时间升序，优先先到期的）
func (r *MerchantRepository) ListAutoRenewEntries(sessionID uint64) ([]model.SeckillEntry, error) {
	var list []model.SeckillEntry
	err := r.db.Where("session_id = ? AND status = ? AND auto_renew = 1", sessionID, model.SeckillEntryActive).
		Order("id ASC").Find(&list).Error
	return list, err
}

func (r *MerchantRepository) SetSeckillAutoRenew(shopID, entryID uint64, autoRenew int8) error {
	if autoRenew != 0 {
		autoRenew = 1
	}
	res := r.db.Model(&model.SeckillEntry{}).
		Where("id = ? AND shop_id = ?", entryID, shopID).
		Update("auto_renew", autoRenew)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("报名记录不存在")
	}
	return nil
}

func (r *MerchantRepository) FindShopEntry(shopID, entryID uint64) (*model.SeckillEntry, error) {
	var e model.SeckillEntry
	if err := r.db.Where("id = ? AND shop_id = ?", entryID, shopID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

// RenewSeckillEntry 将旧报名续到新场次（余额不足则返回错误，由调用方跳过）
func (r *MerchantRepository) RenewSeckillEntry(old *model.SeckillEntry, newSessionID uint64, fee float64, maxPerShop int) error {
	cnt, err := r.CountShopEntries(newSessionID, old.ShopID)
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
	return r.applySeckillEntryTx(entry, fee, nil, model.WalletLogSeckillRenew, "秒杀到期自动续费")
}

func isDup(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate") || strings.Contains(msg, "UNIQUE")
}

func (r *MerchantRepository) ListAdminEntries(sessionID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	q := r.db.Model(&model.SeckillEntry{})
	if sessionID > 0 {
		q = q.Where("session_id = ?", sessionID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SeckillEntry
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) ListShopEntries(shopID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	q := r.db.Model(&model.SeckillEntry{}).Where("shop_id = ?", shopID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SeckillEntry
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) ListActiveEntries(sessionID uint64, limit int) ([]model.SeckillEntry, error) {
	var list []model.SeckillEntry
	err := r.db.Where("session_id = ? AND status = ?", sessionID, model.SeckillEntryActive).
		Order("id ASC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *MerchantRepository) ListActiveEntriesPage(sessionID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	q := r.db.Model(&model.SeckillEntry{}).Where("session_id = ? AND status = ?", sessionID, model.SeckillEntryActive)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SeckillEntry
	err := q.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MerchantRepository) FindSeckillEntry(id uint64) (*model.SeckillEntry, error) {
	var e model.SeckillEntry
	if err := r.db.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *MerchantRepository) DecrSeckillStock(entryID uint64, qty int) error {
	res := r.db.Exec(
		`UPDATE seckill_entries SET seckill_stock = seckill_stock - ? WHERE id = ? AND status = ? AND seckill_stock >= ?`,
		qty, entryID, model.SeckillEntryActive, qty,
	)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("秒杀库存不足或场次无效")
	}
	return nil
}

func (r *MerchantRepository) IncrSeckillStock(entryID uint64, qty int) error {
	if qty <= 0 {
		return nil
	}
	return r.db.Exec(
		`UPDATE seckill_entries SET seckill_stock = seckill_stock + ? WHERE id = ?`,
		qty, entryID,
	).Error
}
