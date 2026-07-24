package biz

import (
	"context"
	"errors"
	"strings"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/types"
)

func validSlotType(t string) bool {
	switch t {
	case model.SlotBrandShop, model.SlotQualityShop, model.SlotArticle:
		return true
	}
	return false
}

func (l *MerchantLogic) ListSlotPackages(slotType string, onlyOn bool) ([]model.HomepageSlotPackage, error) {
	return l.svcCtx.Repo.ListSlotPackages(context.Background(), slotType, onlyOn)
}

func (l *MerchantLogic) CreateSlotPackage(p *model.HomepageSlotPackage) error {
	if !validSlotType(p.SlotType) {
		return errors.New("无效展位类型")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("套餐名称不能为空")
	}
	if p.Price < 0 || p.DurationDays < 1 {
		return errors.New("价格或时长无效")
	}
	if p.Status == "" {
		p.Status = model.SlotPkgOn
	}
	return l.svcCtx.Repo.CreateSlotPackage(context.Background(), p)
}

func (l *MerchantLogic) UpdateSlotPackage(id uint64, p *model.HomepageSlotPackage) error {
	old, err := l.svcCtx.Repo.GetSlotPackage(context.Background(), id)
	if err != nil {
		return errors.New("套餐不存在")
	}
	p.ID = old.ID
	if !validSlotType(p.SlotType) {
		return errors.New("无效展位类型")
	}
	return l.svcCtx.Repo.UpdateSlotPackage(context.Background(), p)
}

func (l *MerchantLogic) ListSlotSettings() ([]model.HomepageSlotSetting, error) {
	list, err := l.svcCtx.Repo.ListSlotSettings(context.Background())
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		for _, t := range []string{model.SlotBrandShop, model.SlotQualityShop, model.SlotArticle} {
			_, _ = l.svcCtx.Repo.GetSlotSetting(context.Background(), t)
		}
		return l.svcCtx.Repo.ListSlotSettings(context.Background())
	}
	return list, nil
}

func (l *MerchantLogic) UpdateSlotSettings(items []model.HomepageSlotSetting) error {
	for _, it := range items {
		if !validSlotType(it.SlotType) {
			continue
		}
		if it.HomeLimit < 1 {
			it.HomeLimit = 1
		}
		if err := l.svcCtx.Repo.UpsertSlotSetting(context.Background(), it.SlotType, it.HomeLimit); err != nil {
			return err
		}
	}
	return nil
}

func (l *MerchantLogic) ListSlotOrders(shopID uint64, slotType, status string, page, pageSize int) ([]model.HomepageSlotOrder, int64, error) {
	l.svcCtx.Repo.ExpireDueSlotOrders(context.Background())
	list, total, err := l.svcCtx.Repo.ListSlotOrders(context.Background(), shopID, slotType, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		if shop, e := l.svcCtx.Repo.FindShop(context.Background(), list[i].ShopID); e == nil {
			list[i].ShopName = shop.Name
		}
		if pkg, e := l.svcCtx.Repo.GetSlotPackage(context.Background(), list[i].PackageID); e == nil {
			list[i].PackageName = pkg.Name
		}
		if list[i].SlotType == model.SlotArticle && list[i].TargetID > 0 {
			if title, e := l.svcCtx.Repo.GetArticleTitle(context.Background(), list[i].TargetID); e == nil {
				list[i].TargetName = title
			}
		} else if list[i].ShopName != "" {
			list[i].TargetName = list[i].ShopName
		}
	}
	return list, total, nil
}

func (l *MerchantLogic) BuySlot(shopID, userID uint64, req types.BuySlotReq) (*model.HomepageSlotOrder, error) {
	pkg, err := l.svcCtx.Repo.GetSlotPackage(context.Background(), req.PackageID)
	if err != nil || pkg.Status != model.SlotPkgOn {
		return nil, errors.New("套餐不存在或已下架")
	}
	targetID := req.TargetID
	if pkg.SlotType == model.SlotArticle {
		if targetID == 0 {
			return nil, errors.New("请选择文章")
		}
		if !l.svcCtx.Repo.ArticlePublishedForShop(context.Background(), targetID, shopID) {
			return nil, errors.New("文章不存在或不属于本店/未发布")
		}
	} else {
		targetID = shopID
	}
	shop, err := l.svcCtx.Repo.FindShop(context.Background(), shopID)
	if err != nil || shop.Status != model.ShopApproved {
		return nil, errors.New("店铺不可用")
	}
	order := &model.HomepageSlotOrder{
		ShopID:       shopID,
		SlotType:     pkg.SlotType,
		PackageID:    pkg.ID,
		TargetID:     targetID,
		Amount:       pkg.Price,
		DurationDays: pkg.DurationDays,
		PaySource:    model.SlotPayWallet,
		OperatorID:   userID,
	}
	op := userID
	if err := l.svcCtx.Repo.PurchaseSlotOrder(context.Background(), order, true, &op); err != nil {
		return nil, err
	}
	return order, nil
}

func (l *MerchantLogic) GrantSlot(adminID uint64, req types.GrantSlotReq) (*model.HomepageSlotOrder, error) {
	pkg, err := l.svcCtx.Repo.GetSlotPackage(context.Background(), req.PackageID)
	if err != nil {
		return nil, errors.New("套餐不存在")
	}
	shop, err := l.svcCtx.Repo.FindShop(context.Background(), req.ShopID)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}
	_ = shop
	targetID := req.TargetID
	if pkg.SlotType == model.SlotArticle {
		if targetID == 0 {
			return nil, errors.New("请指定文章")
		}
	} else {
		targetID = req.ShopID
	}
	order := &model.HomepageSlotOrder{
		ShopID:       req.ShopID,
		SlotType:     pkg.SlotType,
		PackageID:    pkg.ID,
		TargetID:     targetID,
		Amount:       0,
		DurationDays: pkg.DurationDays,
		PaySource:    model.SlotPayAdmin,
		OperatorID:   adminID,
	}
	if err := l.svcCtx.Repo.PurchaseSlotOrder(context.Background(), order, false, &adminID); err != nil {
		return nil, err
	}
	return order, nil
}

func (l *MerchantLogic) HomeSlots(slotType string, city string) ([]map[string]interface{}, error) {
	if !validSlotType(slotType) || slotType == model.SlotArticle {
		return nil, errors.New("无效展位类型")
	}
	setting, _ := l.svcCtx.Repo.GetSlotSetting(context.Background(), slotType)
	limit := 8
	if setting != nil && setting.HomeLimit > 0 {
		limit = setting.HomeLimit
	}
	list, _, err := l.svcCtx.Repo.ListPublicShopsWithSlot(context.Background(), slotType, 1, limit, limit, city)
	if err != nil {
		return nil, err
	}
	paid, _ := l.svcCtx.Repo.ActivePaidTargetIDs(context.Background(), slotType)
	out := make([]map[string]interface{}, 0, len(list))
	for _, s := range list {
		out = append(out, map[string]interface{}{
			"id": s.ID, "name": s.Name, "logo": s.Logo, "category": s.Category,
			"storefront_image": s.StorefrontImage, "description": s.Description,
			"city": s.City, "paid": paid[s.ID],
		})
	}
	return out, nil
}

func (l *MerchantLogic) ListPublicShopsSlot(slotType string, page, pageSize int, city string) ([]map[string]interface{}, int64, error) {
	if slotType == "" {
		slotType = model.SlotQualityShop
	}
	if !validSlotType(slotType) || slotType == model.SlotArticle {
		return nil, 0, errors.New("无效展位类型")
	}
	list, total, err := l.svcCtx.Repo.ListPublicShopsWithSlot(context.Background(), slotType, page, pageSize, 0, city)
	if err != nil {
		return nil, 0, err
	}
	paid, _ := l.svcCtx.Repo.ActivePaidTargetIDs(context.Background(), slotType)
	out := make([]map[string]interface{}, 0, len(list))
	for _, s := range list {
		out = append(out, map[string]interface{}{
			"id": s.ID, "name": s.Name, "logo": s.Logo, "category": s.Category,
			"storefront_image": s.StorefrontImage, "description": s.Description,
			"city": s.City, "paid": paid[s.ID],
		})
	}
	return out, total, nil
}
