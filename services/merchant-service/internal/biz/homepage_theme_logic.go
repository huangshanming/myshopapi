package biz

import (
	"context"
	"errors"
	"strings"

	"mymall/services/merchant-service/internal/model"
)

type ThemeBuyReq struct {
	ThemeSlotID uint64 `json:"theme_slot_id"`
	PackageID   uint64 `json:"package_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	CoverURL    string `json:"cover_url"`
	LinkType    string `json:"link_type"`
	LinkID      uint64 `json:"link_id"`
}

type ThemeGrantReq struct {
	ShopID      uint64 `json:"shop_id"`
	ThemeSlotID uint64 `json:"theme_slot_id"`
	PackageID   uint64 `json:"package_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	CoverURL    string `json:"cover_url"`
	LinkType    string `json:"link_type"`
	LinkID      uint64 `json:"link_id"`
}

func (l *MerchantLogic) validateThemeCreative(shopID uint64, title, cover, linkType string, linkID uint64) (string, uint64, error) {
	title = strings.TrimSpace(title)
	cover = strings.TrimSpace(cover)
	if title == "" {
		return "", 0, errors.New("请填写标题")
	}
	if cover == "" {
		return "", 0, errors.New("请上传封面图")
	}
	lt := strings.TrimSpace(linkType)
	switch lt {
	case model.ThemeLinkShop:
		if linkID == 0 {
			linkID = shopID
		}
		if linkID != shopID {
			return "", 0, errors.New("跳转店铺必须为本店")
		}
	case model.ThemeLinkCategory:
		if linkID == 0 || !l.svcCtx.Repo.CategoryExistsShow(context.Background(), linkID) {
			return "", 0, errors.New("请选择有效分类")
		}
	case model.ThemeLinkProduct:
		if linkID == 0 || !l.svcCtx.Repo.ProductOnSaleOfShop(context.Background(), linkID, shopID) {
			return "", 0, errors.New("请选择本店在售商品")
		}
	default:
		return "", 0, errors.New("跳转类型无效")
	}
	return lt, linkID, nil
}

func (l *MerchantLogic) ListThemeTiles() ([]model.ThemeTile, error) {
	return l.svcCtx.Repo.BuildThemeTiles(context.Background())
}

func (l *MerchantLogic) AdminListThemeSlots() ([]model.HomepageThemeSlot, error) {
	return l.svcCtx.Repo.ListThemeSlots(context.Background(), false)
}

func (l *MerchantLogic) AdminUpdateThemeSlot(id uint64, updates map[string]interface{}) error {
	return l.svcCtx.Repo.UpdateThemeSlot(context.Background(), id, updates)
}

func (l *MerchantLogic) ListThemePackages(themeSlotID uint64, onlyOn bool) ([]model.HomepageThemePackage, error) {
	return l.svcCtx.Repo.ListThemePackages(context.Background(), themeSlotID, onlyOn)
}

func (l *MerchantLogic) AdminCreateThemePackage(p *model.HomepageThemePackage) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("名称不能为空")
	}
	if p.DurationDays < 1 {
		return errors.New("天数无效")
	}
	if p.Status == "" {
		p.Status = model.ThemeSlotOn
	}
	return l.svcCtx.Repo.CreateThemePackage(context.Background(), p)
}

func (l *MerchantLogic) AdminUpdateThemePackage(id uint64, updates map[string]interface{}) error {
	return l.svcCtx.Repo.UpdateThemePackage(context.Background(), id, updates)
}

func (l *MerchantLogic) ListThemeOrders(shopID, themeSlotID uint64, page, pageSize int) ([]model.HomepageThemeOrder, int64, error) {
	list, total, err := l.svcCtx.Repo.ListThemeOrders(context.Background(), shopID, themeSlotID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		if shop, e := l.svcCtx.Repo.FindShop(context.Background(), list[i].ShopID); e == nil {
			list[i].ShopName = shop.Name
		}
		if s, e := l.svcCtx.Repo.GetThemeSlot(context.Background(), list[i].ThemeSlotID); e == nil {
			list[i].ThemeSlotName = s.Name
		}
		if p, e := l.svcCtx.Repo.GetThemePackage(context.Background(), list[i].PackageID); e == nil {
			list[i].PackageName = p.Name
		}
	}
	return list, total, nil
}

func (l *MerchantLogic) BuyTheme(shopID, userID uint64, req ThemeBuyReq) (*model.HomepageThemeOrder, error) {
	slot, err := l.svcCtx.Repo.GetThemeSlot(context.Background(), req.ThemeSlotID)
	if err != nil || slot.Status != model.ThemeSlotOn {
		return nil, errors.New("坑位不存在或已下架")
	}
	pkg, err := l.svcCtx.Repo.GetThemePackage(context.Background(), req.PackageID)
	if err != nil || pkg.Status != model.ThemeSlotOn {
		return nil, errors.New("套餐不存在或已下架")
	}
	if pkg.ThemeSlotID > 0 && pkg.ThemeSlotID != req.ThemeSlotID {
		return nil, errors.New("套餐不适用于该坑位")
	}
	shop, err := l.svcCtx.Repo.FindShop(context.Background(), shopID)
	if err != nil || shop.Status != model.ShopApproved {
		return nil, errors.New("店铺不可用")
	}
	lt, lid, err := l.validateThemeCreative(shopID, req.Title, req.CoverURL, req.LinkType, req.LinkID)
	if err != nil {
		return nil, err
	}
	cover := strings.TrimSpace(req.CoverURL)
	if cover == "" {
		cover = shop.StorefrontImage
		if cover == "" {
			cover = shop.Logo
		}
	}
	order := &model.HomepageThemeOrder{
		ShopID:       shopID,
		ThemeSlotID:  req.ThemeSlotID,
		PackageID:    pkg.ID,
		Title:        strings.TrimSpace(req.Title),
		Subtitle:     strings.TrimSpace(req.Subtitle),
		CoverURL:     cover,
		LinkType:     lt,
		LinkID:       lid,
		Amount:       pkg.Price,
		DurationDays: pkg.DurationDays,
		PaySource:    model.ThemePayWallet,
		OperatorID:   userID,
	}
	uid := userID
	if err := l.svcCtx.Repo.PurchaseThemeOrder(context.Background(), order, true, &uid); err != nil {
		return nil, err
	}
	return order, nil
}

func (l *MerchantLogic) GrantTheme(adminID uint64, req ThemeGrantReq) (*model.HomepageThemeOrder, error) {
	if req.ShopID == 0 {
		return nil, errors.New("请选择店铺")
	}
	slot, err := l.svcCtx.Repo.GetThemeSlot(context.Background(), req.ThemeSlotID)
	if err != nil || slot.Status != model.ThemeSlotOn {
		return nil, errors.New("坑位不存在或已下架")
	}
	pkg, err := l.svcCtx.Repo.GetThemePackage(context.Background(), req.PackageID)
	if err != nil {
		return nil, errors.New("套餐不存在")
	}
	if pkg.ThemeSlotID > 0 && pkg.ThemeSlotID != req.ThemeSlotID {
		return nil, errors.New("套餐不适用于该坑位")
	}
	shop, err := l.svcCtx.Repo.FindShop(context.Background(), req.ShopID)
	if err != nil || shop.Status != model.ShopApproved {
		return nil, errors.New("店铺不可用")
	}
	lt, lid, err := l.validateThemeCreative(req.ShopID, req.Title, req.CoverURL, req.LinkType, req.LinkID)
	if err != nil {
		return nil, err
	}
	cover := strings.TrimSpace(req.CoverURL)
	if cover == "" {
		cover = shop.StorefrontImage
		if cover == "" {
			cover = shop.Logo
		}
	}
	order := &model.HomepageThemeOrder{
		ShopID:       req.ShopID,
		ThemeSlotID:  req.ThemeSlotID,
		PackageID:    pkg.ID,
		Title:        strings.TrimSpace(req.Title),
		Subtitle:     strings.TrimSpace(req.Subtitle),
		CoverURL:     cover,
		LinkType:     lt,
		LinkID:       lid,
		Amount:       0,
		DurationDays: pkg.DurationDays,
		PaySource:    model.ThemePayAdmin,
		OperatorID:   adminID,
	}
	oid := adminID
	if err := l.svcCtx.Repo.PurchaseThemeOrder(context.Background(), order, false, &oid); err != nil {
		return nil, err
	}
	return order, nil
}
