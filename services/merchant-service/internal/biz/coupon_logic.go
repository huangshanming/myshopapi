package biz

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/types"

	"github.com/google/uuid"
)

type MatchCouponView struct {
	UserCouponID    uint64  `json:"user_coupon_id"`
	CouponID        uint64  `json:"coupon_id"`
	Name            string  `json:"name"`
	CouponType      string  `json:"coupon_type"`
	DiscountAmount  float64 `json:"discount_amount"`
	ThresholdAmount float64 `json:"threshold_amount"`
	ValidEnd        string  `json:"valid_end"`
	Usable          bool    `json:"usable"`
	Reason          string  `json:"reason,omitempty"`
	Best            bool    `json:"best,omitempty"`
}

type MatchResp struct {
	GoodsAmount      float64           `json:"goods_amount"`
	DiscountAmount   float64           `json:"discount_amount"`
	PayAmount        float64           `json:"pay_amount"`
	BestUserCouponID uint64            `json:"best_user_coupon_id"`
	Available        []MatchCouponView `json:"available"`
	Unavailable      []MatchCouponView `json:"unavailable"`
}

func (l *MerchantLogic) validateCouponSave(req types.CouponSaveReq, issuerType string, shopID uint64) (*model.Coupon, []model.CouponScope, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, nil, errors.New("请填写名称")
	}
	ct := req.CouponType
	scopeType := req.ScopeType
	if scopeType == "" {
		scopeType = model.CouponScopeAll
	}
	switch ct {
	case model.CouponTypeFullReduce:
		if req.ThresholdAmount <= 0 || req.DiscountAmount <= 0 {
			return nil, nil, errors.New("满减券需设置门槛与减免额")
		}
		scopeType = model.CouponScopeAll
	case model.CouponTypeNoThreshold:
		req.ThresholdAmount = 0
		if req.DiscountAmount <= 0 {
			return nil, nil, errors.New("无门槛券需设置减免额")
		}
		scopeType = model.CouponScopeAll
	case model.CouponTypeCategory:
		scopeType = model.CouponScopeCategory
		if req.DiscountAmount <= 0 {
			return nil, nil, errors.New("品类券需设置减免额")
		}
		if len(req.Scopes) == 0 {
			return nil, nil, errors.New("请选择适用分类")
		}
	case model.CouponTypeProduct:
		scopeType = model.CouponScopeProduct
		if req.DiscountAmount <= 0 {
			return nil, nil, errors.New("商品券需设置减免额")
		}
		if len(req.Scopes) == 0 {
			return nil, nil, errors.New("请选择适用商品")
		}
	case model.CouponTypeDiscount:
		if req.DiscountRate <= 0 || req.DiscountRate >= 1 {
			return nil, nil, errors.New("折扣比例需在 0~1 之间，如 0.8 表示八折")
		}
		if scopeType == "" {
			scopeType = model.CouponScopeAll
		}
	default:
		return nil, nil, errors.New("券类型无效")
	}
	if req.PerUserLimit < 1 {
		return nil, nil, errors.New("每人限领至少 1 张")
	}
	if req.TotalCount < 0 {
		return nil, nil, errors.New("发放总量无效")
	}
	vt := req.ValidType
	if vt == "" {
		vt = model.CouponValidFixed
	}
	if vt == model.CouponValidFixed {
		if req.ValidStart == nil || req.ValidEnd == nil {
			return nil, nil, errors.New("请设置固定有效期")
		}
		if time.Time(*req.ValidEnd).Before(time.Time(*req.ValidStart)) {
			return nil, nil, errors.New("结束时间需晚于开始时间")
		}
	} else if vt == model.CouponValidRelative {
		if req.ValidDays < 1 {
			return nil, nil, errors.New("领取后有效天数至少 1 天")
		}
	} else {
		return nil, nil, errors.New("有效期类型无效")
	}
	channels := req.Channels
	if len(channels) == 0 {
		channels = []string{model.CouponChannelDirect}
	}
	identity := req.UserIdentity
	if identity == "" {
		identity = model.CouponIdentityAll
	}
	status := req.Status
	if status == "" {
		status = model.CouponStatusOn
	}
	stackable := int8(0)
	if req.Stackable != nil {
		stackable = *req.Stackable
	}
	scopes := make([]model.CouponScope, 0, len(req.Scopes))
	for _, s := range req.Scopes {
		if s.RefID == 0 {
			continue
		}
		rt := s.RefType
		if rt == "" {
			rt = scopeType
		}
		if scopeType == model.CouponScopeCategory {
			rt = model.CouponScopeCategory
			if !l.svcCtx.Repo.CategoryExistsShow(context.Background(), s.RefID) {
				return nil, nil, fmt.Errorf("分类 %d 无效", s.RefID)
			}
		}
		if scopeType == model.CouponScopeProduct {
			rt = model.CouponScopeProduct
			if issuerType == model.CouponIssuerShop {
				if !l.svcCtx.Repo.ProductOnSaleOfShop(context.Background(), s.RefID, shopID) {
					return nil, nil, fmt.Errorf("商品 %d 非本店在售", s.RefID)
				}
			}
		}
		scopes = append(scopes, model.CouponScope{RefType: rt, RefID: s.RefID})
	}
	if (ct == model.CouponTypeCategory || ct == model.CouponTypeProduct) && len(scopes) == 0 {
		return nil, nil, errors.New("适用范围不能为空")
	}
	c := &model.Coupon{
		Name:              name,
		IssuerType:        issuerType,
		ShopID:            shopID,
		CouponType:        ct,
		ThresholdAmount:   req.ThresholdAmount,
		DiscountAmount:    req.DiscountAmount,
		DiscountRate:      req.DiscountRate,
		MaxDiscountAmount: req.MaxDiscountAmount,
		ScopeType:         scopeType,
		TotalCount:        req.TotalCount,
		PerUserLimit:      req.PerUserLimit,
		ValidType:         vt,
		ValidStart:        req.ValidStart,
		ValidEnd:          req.ValidEnd,
		ValidDays:         req.ValidDays,
		Stackable:         stackable,
		UserIdentity:      identity,
		Channels:          model.StringSlice(channels),
		Status:            status,
		Remark:            strings.TrimSpace(req.Remark),
	}
	return c, scopes, nil
}

func (l *MerchantLogic) AdminCreateCoupon(adminID uint64, req types.CouponSaveReq) (*model.Coupon, error) {
	c, scopes, err := l.validateCouponSave(req, model.CouponIssuerPlatform, 0)
	if err != nil {
		return nil, err
	}
	c.CreatedBy = adminID
	if err := l.svcCtx.Repo.CreateCoupon(context.Background(), c, scopes); err != nil {
		return nil, err
	}
	c.Scopes = scopes
	return c, nil
}

func (l *MerchantLogic) MerchantCreateCoupon(shopID, userID uint64, req types.CouponSaveReq) (*model.Coupon, error) {
	if shopID == 0 {
		return nil, errors.New("缺少店铺")
	}
	c, scopes, err := l.validateCouponSave(req, model.CouponIssuerShop, shopID)
	if err != nil {
		return nil, err
	}
	c.CreatedBy = userID
	if err := l.svcCtx.Repo.CreateCoupon(context.Background(), c, scopes); err != nil {
		return nil, err
	}
	c.Scopes = scopes
	return c, nil
}

func (l *MerchantLogic) UpdateCoupon(id, shopID uint64, platform bool, req types.CouponSaveReq) error {
	old, err := l.svcCtx.Repo.GetCoupon(context.Background(), id)
	if err != nil {
		return errors.New("优惠券不存在")
	}
	if platform {
		if old.IssuerType != model.CouponIssuerPlatform {
			return errors.New("无权编辑")
		}
	} else {
		if old.IssuerType != model.CouponIssuerShop || old.ShopID != shopID {
			return errors.New("无权编辑")
		}
	}
	c, scopes, err := l.validateCouponSave(req, old.IssuerType, old.ShopID)
	if err != nil {
		return err
	}
	// 发放中限制：仅允许提高总量、下架等；此处若已有领取，禁止改类型与面额
	if old.ClaimedCount > 0 {
		updates := map[string]interface{}{
			"name":        c.Name,
			"total_count": c.TotalCount,
			"status":      c.Status,
			"remark":      c.Remark,
			"channels":    c.Channels,
			"stackable":   c.Stackable,
		}
		if c.TotalCount > 0 && c.TotalCount < old.ClaimedCount {
			return errors.New("发放总量不能小于已领数量")
		}
		return l.svcCtx.Repo.UpdateCoupon(context.Background(), id, updates, nil)
	}
	updates := map[string]interface{}{
		"name": c.Name, "coupon_type": c.CouponType, "threshold_amount": c.ThresholdAmount,
		"discount_amount": c.DiscountAmount, "discount_rate": c.DiscountRate,
		"max_discount_amount": c.MaxDiscountAmount, "scope_type": c.ScopeType,
		"total_count": c.TotalCount, "per_user_limit": c.PerUserLimit,
		"valid_type": c.ValidType, "valid_start": c.ValidStart, "valid_end": c.ValidEnd,
		"valid_days": c.ValidDays, "stackable": c.Stackable, "user_identity": c.UserIdentity,
		"channels": c.Channels, "status": c.Status, "remark": c.Remark,
	}
	return l.svcCtx.Repo.UpdateCoupon(context.Background(), id, updates, &scopes)
}

func (l *MerchantLogic) OffCoupon(id, shopID uint64, platform bool) error {
	old, err := l.svcCtx.Repo.GetCoupon(context.Background(), id)
	if err != nil {
		return errors.New("优惠券不存在")
	}
	if platform {
		if old.IssuerType != model.CouponIssuerPlatform {
			return errors.New("无权操作")
		}
	} else if old.ShopID != shopID {
		return errors.New("无权操作")
	}
	return l.svcCtx.Repo.UpdateCoupon(context.Background(), id, map[string]interface{}{"status": model.CouponStatusOff}, nil)
}

func (l *MerchantLogic) CopyCoupon(id, shopID, operatorID uint64, platform bool) (*model.Coupon, error) {
	old, err := l.svcCtx.Repo.GetCoupon(context.Background(), id)
	if err != nil {
		return nil, errors.New("优惠券不存在")
	}
	if platform {
		if old.IssuerType != model.CouponIssuerPlatform {
			return nil, errors.New("无权操作")
		}
	} else if old.ShopID != shopID {
		return nil, errors.New("无权操作")
	}
	nc := *old
	nc.ID = 0
	nc.Name = old.Name + "（副本）"
	nc.ClaimedCount = 0
	nc.Status = model.CouponStatusDraft
	nc.CreatedBy = operatorID
	scopes := append([]model.CouponScope{}, old.Scopes...)
	if err := l.svcCtx.Repo.CreateCoupon(context.Background(), &nc, scopes); err != nil {
		return nil, err
	}
	nc.Scopes = scopes
	return &nc, nil
}

func (l *MerchantLogic) ListCoupons(issuerType string, shopID uint64, status, keyword string, page, pageSize int) ([]model.Coupon, int64, error) {
	return l.svcCtx.Repo.ListCoupons(context.Background(), issuerType, shopID, status, keyword, page, pageSize)
}

func (l *MerchantLogic) GetCoupon(id uint64) (*model.Coupon, error) {
	return l.svcCtx.Repo.GetCoupon(context.Background(), id)
}

func channelAllowed(channels model.StringSlice, source string) bool {
	if len(channels) == 0 {
		return source == model.CouponSourceDirect || source == model.CouponSourcePopup
	}
	for _, ch := range channels {
		if ch == source || (source == model.CouponSourceDirect && ch == model.CouponChannelDirect) ||
			(source == model.CouponSourcePopup && ch == model.CouponChannelPopup) ||
			(source == model.CouponSourceTargeted && ch == model.CouponChannelTargeted) ||
			(source == model.CouponSourceOrderGift && ch == model.CouponChannelOrderGift) {
			return true
		}
	}
	return false
}

func (l *MerchantLogic) checkIdentity(userID uint64, identity string) error {
	if identity == "" || identity == model.CouponIdentityAll {
		return nil
	}
	created, err := l.svcCtx.Repo.UserCreatedAt(context.Background(), userID)
	if err != nil || created.IsZero() {
		return errors.New("用户无效")
	}
	isNew := time.Since(created) <= 7*24*time.Hour
	if identity == model.CouponIdentityNew && !isNew {
		return errors.New("仅新用户可领")
	}
	if identity == model.CouponIdentityOld && isNew {
		return errors.New("仅老用户可领")
	}
	return nil
}

func (l *MerchantLogic) ClaimCoupon(userID, couponID uint64, source string) (*model.UserCoupon, error) {
	if userID == 0 {
		return nil, errors.New("请先登录")
	}
	if source == "" {
		source = model.CouponSourceDirect
	}
	c, err := l.svcCtx.Repo.GetCoupon(context.Background(), couponID)
	if err != nil {
		return nil, errors.New("优惠券不存在")
	}
	if c.Status != model.CouponStatusOn {
		return nil, errors.New("优惠券不可领取")
	}
	if !channelAllowed(c.Channels, source) {
		return nil, errors.New("该渠道不可领取")
	}
	if err := l.checkIdentity(userID, c.UserIdentity); err != nil {
		return nil, err
	}
	return l.svcCtx.Repo.ClaimCoupon(context.Background(), userID, c, source, "")
}

func (l *MerchantLogic) GrantCoupon(operatorID uint64, couponID uint64, userIDs []uint64, shopID uint64, platform bool) (*model.CouponGrant, error) {
	c, err := l.svcCtx.Repo.GetCoupon(context.Background(), couponID)
	if err != nil {
		return nil, errors.New("优惠券不存在")
	}
	if platform {
		if c.IssuerType != model.CouponIssuerPlatform {
			return nil, errors.New("无权发放")
		}
	} else if c.ShopID != shopID {
		return nil, errors.New("无权发放")
	}
	if !channelAllowed(c.Channels, model.CouponSourceTargeted) && len(c.Channels) > 0 {
		// 允许 targeted 即使未配置：运营代发
	}
	batch := uuid.NewString()[:8]
	ok := 0
	for _, uid := range userIDs {
		if uid == 0 {
			continue
		}
		if _, err := l.svcCtx.Repo.ClaimCoupon(context.Background(), uid, c, model.CouponSourceTargeted, batch); err == nil {
			ok++
		}
	}
	g := &model.CouponGrant{
		CouponID:     couponID,
		OperatorID:   operatorID,
		IssuerType:   c.IssuerType,
		ShopID:       c.ShopID,
		UserCount:    len(userIDs),
		SuccessCount: ok,
		BatchNo:      batch,
	}
	_ = l.svcCtx.Repo.CreateGrant(context.Background(), g)
	return g, nil
}

func (l *MerchantLogic) ListCenter(userID, shopID uint64) ([]model.Coupon, error) {
	list, err := l.svcCtx.Repo.ListCenterCoupons(context.Background(), shopID)
	if err != nil {
		return nil, err
	}
	if userID > 0 {
		for i := range list {
			n, _ := l.svcCtx.Repo.CountUserClaims(context.Background(), list[i].ID, userID)
			list[i].ClaimedByMe = n > 0 && int(n) >= list[i].PerUserLimit
		}
	}
	return list, nil
}

func (l *MerchantLogic) ListPopup(userID uint64) ([]model.Coupon, error) {
	list, err := l.svcCtx.Repo.ListPopupCoupons(context.Background())
	if err != nil {
		return nil, err
	}
	if userID > 0 {
		for i := range list {
			n, _ := l.svcCtx.Repo.CountUserClaims(context.Background(), list[i].ID, userID)
			list[i].ClaimedByMe = int(n) >= list[i].PerUserLimit
		}
	}
	return list, nil
}

func (l *MerchantLogic) ListMyCoupons(userID uint64, status string, page, pageSize int) ([]model.UserCoupon, int64, error) {
	return l.svcCtx.Repo.ListUserCoupons(context.Background(), userID, status, page, pageSize)
}

func (l *MerchantLogic) CouponClaims(couponID uint64, page, pageSize int) ([]model.UserCoupon, int64, error) {
	return l.svcCtx.Repo.ListClaims(context.Background(), couponID, page, pageSize)
}

func (l *MerchantLogic) CouponRedeems(couponID uint64, page, pageSize int) ([]model.CouponRedeemLog, int64, error) {
	return l.svcCtx.Repo.ListRedeems(context.Background(), couponID, page, pageSize)
}

func (l *MerchantLogic) CouponStats(couponID uint64) (map[string]interface{}, error) {
	return l.svcCtx.Repo.CouponStats(context.Background(), couponID)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func calcDiscount(c *model.Coupon, eligible float64) float64 {
	if eligible <= 0 {
		return 0
	}
	var d float64
	switch c.CouponType {
	case model.CouponTypeDiscount:
		d = eligible * (1 - c.DiscountRate)
		if c.MaxDiscountAmount > 0 && d > c.MaxDiscountAmount {
			d = c.MaxDiscountAmount
		}
	default:
		d = c.DiscountAmount
	}
	if d > eligible {
		d = eligible
	}
	return round2(d)
}

func eligibleAmount(c *model.Coupon, items []types.MatchItem) float64 {
	scopeIDs := map[uint64]bool{}
	for _, s := range c.Scopes {
		scopeIDs[s.RefID] = true
	}
	var sum float64
	for _, it := range items {
		switch c.ScopeType {
		case model.CouponScopeCategory:
			if scopeIDs[it.CategoryID] {
				sum += it.Amount
			}
		case model.CouponScopeProduct:
			if scopeIDs[it.ProductID] {
				sum += it.Amount
			}
		default:
			sum += it.Amount
		}
	}
	return round2(sum)
}

func hasSeckill(items []types.MatchItem) bool {
	for _, it := range items {
		if it.SeckillEntryID > 0 {
			return true
		}
	}
	return false
}

func (l *MerchantLogic) MatchCoupons(req types.MatchCouponsReq) (*MatchResp, error) {
	goods := 0.0
	for _, it := range req.Items {
		goods += it.Amount
	}
	goods = round2(goods)
	resp := &MatchResp{GoodsAmount: goods, PayAmount: goods, Available: []MatchCouponView{}, Unavailable: []MatchCouponView{}}
	if req.UserID == 0 {
		return resp, nil
	}
	// 补全 category
	ids := make([]uint64, 0, len(req.Items))
	for _, it := range req.Items {
		ids = append(ids, it.ProductID)
	}
	lite, _ := l.svcCtx.Repo.GetProductsLite(context.Background(), ids)
	for i := range req.Items {
		if req.Items[i].CategoryID == 0 {
			if p, ok := lite[req.Items[i].ProductID]; ok {
				req.Items[i].CategoryID = p.CategoryID
			}
		}
	}
	list, err := l.svcCtx.Repo.ListUserUnusedCoupons(context.Background(), req.UserID)
	if err != nil {
		return nil, err
	}
	seckill := hasSeckill(req.Items)
	var best *MatchCouponView
	for _, uc := range list {
		c := uc.Coupon
		if c == nil {
			continue
		}
		view := MatchCouponView{
			UserCouponID:    uc.ID,
			CouponID:        c.ID,
			Name:            c.Name,
			CouponType:      c.CouponType,
			ThresholdAmount: c.ThresholdAmount,
			ValidEnd:        time.Time(uc.ValidEnd).Format("2006-01-02 15:04:05"),
		}
		if c.IssuerType == model.CouponIssuerShop && c.ShopID != req.ShopID {
			view.Usable = false
			view.Reason = "非本店优惠券"
			resp.Unavailable = append(resp.Unavailable, view)
			continue
		}
		if seckill && c.Stackable == 0 {
			view.Usable = false
			view.Reason = "不可与秒杀叠加"
			resp.Unavailable = append(resp.Unavailable, view)
			continue
		}
		eligible := eligibleAmount(c, req.Items)
		if eligible <= 0 {
			view.Usable = false
			view.Reason = "商品不匹配"
			resp.Unavailable = append(resp.Unavailable, view)
			continue
		}
		th := c.ThresholdAmount
		if c.CouponType == model.CouponTypeNoThreshold {
			th = 0
		}
		if eligible < th {
			view.Usable = false
			view.Reason = "金额不足"
			resp.Unavailable = append(resp.Unavailable, view)
			continue
		}
		d := calcDiscount(c, eligible)
		if d <= 0 {
			view.Usable = false
			view.Reason = "无可用抵扣"
			resp.Unavailable = append(resp.Unavailable, view)
			continue
		}
		view.Usable = true
		view.DiscountAmount = d
		resp.Available = append(resp.Available, view)
		if best == nil || view.DiscountAmount > best.DiscountAmount ||
			(view.DiscountAmount == best.DiscountAmount && view.ValidEnd < best.ValidEnd) {
			cp := view
			best = &cp
		}
	}
	selectedID := req.UserCouponID
	if selectedID == 0 && best != nil {
		selectedID = best.UserCouponID
	}
	if best != nil {
		for i := range resp.Available {
			if resp.Available[i].UserCouponID == best.UserCouponID {
				resp.Available[i].Best = true
			}
		}
		resp.BestUserCouponID = best.UserCouponID
	}
	if selectedID > 0 {
		for _, v := range resp.Available {
			if v.UserCouponID == selectedID {
				resp.DiscountAmount = v.DiscountAmount
				break
			}
		}
	}
	pay := round2(goods - resp.DiscountAmount)
	if pay < 0.01 && goods > 0 {
		pay = 0.01
		resp.DiscountAmount = round2(goods - pay)
	}
	resp.PayAmount = pay
	return resp, nil
}

func (l *MerchantLogic) LockCoupon(userCouponID, userID, orderID uint64, discount float64) error {
	return l.svcCtx.Repo.LockUserCoupon(context.Background(), userCouponID, userID, orderID, discount)
}

func (l *MerchantLogic) UnlockCoupon(userCouponID, orderID uint64) error {
	return l.svcCtx.Repo.UnlockUserCoupon(context.Background(), userCouponID, orderID, model.CouponActionUnlock)
}

func (l *MerchantLogic) RedeemCoupon(userCouponID, orderID uint64, discount float64) error {
	return l.svcCtx.Repo.RedeemUserCoupon(context.Background(), userCouponID, orderID, discount)
}

func (l *MerchantLogic) ReturnCoupon(userCouponID, orderID uint64) error {
	return l.svcCtx.Repo.UnlockUserCoupon(context.Background(), userCouponID, orderID, model.CouponActionReturn)
}

func (l *MerchantLogic) OrderGiftCoupons(userID, shopID uint64) (int, error) {
	list, err := l.svcCtx.Repo.ListOrderGiftCoupons(context.Background(), shopID)
	if err != nil {
		return 0, err
	}
	n := 0
	for i := range list {
		if _, err := l.svcCtx.Repo.ClaimCoupon(context.Background(), userID, &list[i], model.CouponSourceOrderGift, ""); err == nil {
			n++
		}
	}
	return n, nil
}
