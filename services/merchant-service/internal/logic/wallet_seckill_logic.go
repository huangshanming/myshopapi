package logic

import (
	"errors"
	"time"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/types"
)

func (l *MerchantLogic) GetWallet(shopID uint64) (*model.ShopWallet, error) {
	if shopID == 0 {
		return nil, errors.New("缺少店铺")
	}
	return l.svcCtx.Repo.GetWallet(shopID)
}

func (l *MerchantLogic) AdjustWallet(shopID uint64, field string, amount float64, remark string, operatorID uint64) (*model.ShopWallet, error) {
	if shopID == 0 {
		return nil, errors.New("店铺无效")
	}
	if amount == 0 {
		return nil, errors.New("调整金额不能为0")
	}
	if field == "" {
		field = model.WalletFieldBalance
	}
	switch field {
	case model.WalletFieldBalance, model.WalletFieldDeposit, model.WalletFieldFrozen:
	default:
		return nil, errors.New("调账类型无效")
	}
	if remark == "" {
		switch field {
		case model.WalletFieldDeposit:
			remark = "平台调整保证金"
		case model.WalletFieldFrozen:
			remark = "平台调整冻结余额"
		default:
			remark = "平台调账"
		}
	}
	op := operatorID
	return l.svcCtx.Repo.AdjustWallet(shopID, field, amount, model.WalletLogAdminAdjust, remark, &op, field, 0)
}

func (l *MerchantLogic) ListWalletLogs(shopID uint64, page, pageSize int) ([]model.ShopWalletLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return l.svcCtx.Repo.ListWalletLogs(shopID, page, pageSize)
}

func (l *MerchantLogic) GetSeckillRule() (*model.SeckillRule, error) {
	return l.svcCtx.Repo.GetActiveSeckillRule()
}

func (l *MerchantLogic) UpdateSeckillRule(req types.SeckillRuleReq) (*model.SeckillRule, error) {
	rule, err := l.svcCtx.Repo.GetActiveSeckillRule()
	if err != nil {
		return nil, err
	}
	if req.DurationHours < 1 {
		return nil, errors.New("场次时长至少1小时")
	}
	if req.MaxEntriesPerShop < 1 {
		return nil, errors.New("每店报名上限至少为1")
	}
	if req.ApplyFee < 0 {
		return nil, errors.New("报名费不能为负")
	}
	rule.DurationHours = req.DurationHours
	rule.ApplyFee = req.ApplyFee
	rule.MaxEntriesPerShop = req.MaxEntriesPerShop
	if req.Status == 0 || req.Status == 1 {
		rule.Status = req.Status
	} else {
		rule.Status = model.SeckillRuleOn
	}
	if err := l.svcCtx.Repo.SaveSeckillRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (l *MerchantLogic) EnsureActiveSession() (*model.SeckillSession, *model.SeckillRule, error) {
	rule, err := l.svcCtx.Repo.GetActiveSeckillRule()
	if err != nil {
		return nil, nil, err
	}
	s, err := l.svcCtx.Repo.GetActiveSession()
	if err == nil {
		return s, rule, nil
	}
	now := time.Now()
	end := now.Add(time.Duration(rule.DurationHours) * time.Hour)
	s, err = l.svcCtx.Repo.CreateSession(rule.ID, now, end)
	return s, rule, err
}

func (l *MerchantLogic) RotateSeckillSessions() {
	rule, err := l.svcCtx.Repo.GetActiveSeckillRule()
	if err != nil || rule.Status != model.SeckillRuleOn {
		return
	}
	s, err := l.svcCtx.Repo.GetActiveSession()
	if err != nil {
		_, _, _ = l.EnsureActiveSession()
		return
	}
	end := time.Time(s.EndAt)
	if !end.IsZero() && !time.Now().Before(end) {
		_ = l.svcCtx.Repo.EndSession(s.ID)
		now := time.Now()
		_, _ = l.svcCtx.Repo.CreateSession(rule.ID, now, now.Add(time.Duration(rule.DurationHours)*time.Hour))
	}
}

func (l *MerchantLogic) ListSeckillSessions(page, pageSize int) ([]model.SeckillSession, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return l.svcCtx.Repo.ListSessions(page, pageSize)
}

func (l *MerchantLogic) ListAdminSeckillEntries(sessionID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return l.svcCtx.Repo.ListAdminEntries(sessionID, page, pageSize)
}

func (l *MerchantLogic) MerchantSeckillSessions() (map[string]interface{}, error) {
	s, rule, err := l.EnsureActiveSession()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"rule": map[string]interface{}{
			"duration_hours":        rule.DurationHours,
			"apply_fee":             rule.ApplyFee,
			"max_entries_per_shop":  rule.MaxEntriesPerShop,
		},
		"sessions": []model.SeckillSession{*s},
	}, nil
}

func (l *MerchantLogic) ApplySeckill(shopID, userID uint64, req types.SeckillApplyReq) (*model.SeckillEntry, error) {
	if shopID == 0 {
		return nil, errors.New("缺少店铺")
	}
	if req.ProductID == 0 || req.ProductName == "" {
		return nil, errors.New("请选择商品")
	}
	if req.SeckillPrice <= 0 {
		return nil, errors.New("秒杀价无效")
	}
	if req.SeckillStock < 1 {
		return nil, errors.New("秒杀库存至少为1")
	}
	s, err := l.svcCtx.Repo.FindSession(req.SessionID)
	if err != nil {
		return nil, errors.New("场次不存在")
	}
	if s.Status != model.SeckillSessionActive {
		return nil, errors.New("场次已结束")
	}
	if !time.Now().Before(time.Time(s.EndAt)) {
		return nil, errors.New("场次已结束")
	}
	rule, err := l.svcCtx.Repo.GetActiveSeckillRule()
	if err != nil {
		return nil, err
	}
	cnt, err := l.svcCtx.Repo.CountShopEntries(s.ID, shopID)
	if err != nil {
		return nil, err
	}
	if int(cnt) >= rule.MaxEntriesPerShop {
		return nil, errors.New("已达本场次报名上限")
	}
	entry := &model.SeckillEntry{
		SessionID:    s.ID,
		ShopID:       shopID,
		ProductID:    req.ProductID,
		ProductName:  req.ProductName,
		ProductImage: req.ProductImage,
		OriginPrice:  req.OriginPrice,
		SeckillPrice: req.SeckillPrice,
		SeckillStock: req.SeckillStock,
	}
	op := userID
	if err := l.svcCtx.Repo.ApplySeckillEntry(entry, rule.ApplyFee, &op); err != nil {
		return nil, err
	}
	return entry, nil
}

func (l *MerchantLogic) ListShopSeckillEntries(shopID uint64, page, pageSize int) ([]model.SeckillEntry, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return l.svcCtx.Repo.ListShopEntries(shopID, page, pageSize)
}

func mapSeckillItem(e model.SeckillEntry) map[string]interface{} {
	return map[string]interface{}{
		"id":         e.ID,
		"product_id": e.ProductID,
		"name":       e.ProductName,
		"price":      e.SeckillPrice,
		"old":        e.OriginPrice,
		"left":       e.SeckillStock,
		"img":        e.ProductImage,
		"shop_id":    e.ShopID,
	}
}

func (l *MerchantLogic) PublicSeckillCurrent() (map[string]interface{}, error) {
	s, _, err := l.EnsureActiveSession()
	if err != nil {
		return nil, err
	}
	entries, err := l.svcCtx.Repo.ListActiveEntries(s.ID, 50)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		items = append(items, mapSeckillItem(e))
	}
	return map[string]interface{}{
		"session_id": s.ID,
		"start_at":   s.StartAt,
		"end_at":     s.EndAt,
		"items":      items,
	}, nil
}

func (l *MerchantLogic) PublicSeckillList(page, pageSize int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}
	s, _, err := l.EnsureActiveSession()
	if err != nil {
		return nil, err
	}
	entries, total, err := l.svcCtx.Repo.ListActiveEntriesPage(s.ID, page, pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		items = append(items, mapSeckillItem(e))
	}
	return map[string]interface{}{
		"session_id": s.ID,
		"start_at":   s.StartAt,
		"end_at":     s.EndAt,
		"total":      total,
		"list":       items,
	}, nil
}

func (l *MerchantLogic) PublicSeckillEntry(id uint64) (map[string]interface{}, error) {
	e, err := l.svcCtx.Repo.FindSeckillEntry(id)
	if err != nil {
		return nil, errors.New("秒杀商品不存在")
	}
	s, err := l.svcCtx.Repo.FindSession(e.SessionID)
	sessionActive := err == nil && s.Status == model.SeckillSessionActive && time.Now().Before(time.Time(s.EndAt))
	available := sessionActive && e.Status == model.SeckillEntryActive && e.SeckillStock > 0
	out := map[string]interface{}{
		"id":              e.ID,
		"session_id":      e.SessionID,
		"product_id":      e.ProductID,
		"product_name":    e.ProductName,
		"product_image":   e.ProductImage,
		"origin_price":    e.OriginPrice,
		"seckill_price":   e.SeckillPrice,
		"seckill_stock":   e.SeckillStock,
		"shop_id":         e.ShopID,
		"status":          e.Status,
		"session_active":  sessionActive,
		"seckill_available": available,
	}
	if err == nil {
		out["start_at"] = s.StartAt
		out["end_at"] = s.EndAt
		out["session_status"] = s.Status
	}
	return out, nil
}

func (l *MerchantLogic) ConsumeSeckill(entryID, productID uint64, qty int) (map[string]interface{}, error) {
	if entryID == 0 || productID == 0 || qty < 1 {
		return nil, errors.New("参数无效")
	}
	e, err := l.svcCtx.Repo.FindSeckillEntry(entryID)
	if err != nil {
		return nil, errors.New("秒杀商品不存在")
	}
	if e.ProductID != productID {
		return nil, errors.New("秒杀商品不匹配")
	}
	if e.Status != model.SeckillEntryActive {
		return nil, errors.New("秒杀报名已失效")
	}
	s, err := l.svcCtx.Repo.FindSession(e.SessionID)
	if err != nil || s.Status != model.SeckillSessionActive || !time.Now().Before(time.Time(s.EndAt)) {
		return nil, errors.New("秒杀场次已结束")
	}
	if err := l.svcCtx.Repo.DecrSeckillStock(entryID, qty); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"entry_id":      entryID,
		"product_id":    productID,
		"seckill_price": e.SeckillPrice,
		"quantity":      qty,
	}, nil
}

func (l *MerchantLogic) RestoreSeckill(entryID uint64, qty int) error {
	if entryID == 0 || qty < 1 {
		return nil
	}
	return l.svcCtx.Repo.IncrSeckillStock(entryID, qty)
}
