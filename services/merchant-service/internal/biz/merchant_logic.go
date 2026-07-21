package biz

import (
	"context"
	"errors"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"gorm.io/gorm"
)

type MerchantLogic struct {
	svcCtx *svc.ServiceContext
}

func NewMerchantLogic(svcCtx *svc.ServiceContext) *MerchantLogic {
	return &MerchantLogic{svcCtx: svcCtx}
}

func (l *MerchantLogic) Apply(ctx context.Context, userID uint64, in types.ApplyReq) (*model.ShopApplication, error) {
	if _, err := l.svcCtx.Repo.FindPendingAppByUser(ctx, userID); err == nil {
		return nil, errors.New("已有待审核的入驻申请")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	app := &model.ShopApplication{
		UserID:            userID,
		ShopName:          in.ShopName,
		ContactName:       in.ContactName,
		ContactPhone:      in.ContactPhone,
		Description:       in.Description,
		Category:          in.Category,
		Province:          in.Province,
		City:              in.City,
		District:          in.District,
		Address:           in.Address,
		BusinessLicenseNo: in.BusinessLicenseNo,
		LegalPerson:       in.LegalPerson,
		LicenseImage:      in.LicenseImage,
		StorefrontImage:   in.StorefrontImage,
		Status:            model.AppPending,
	}
	if err := l.svcCtx.Repo.CreateApplication(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (l *MerchantLogic) ListApplications(ctx context.Context, status string, page, pageSize int) ([]model.ShopApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListApplications(ctx, status, page, pageSize)
}

func (l *MerchantLogic) Approve(ctx context.Context, appID, adminID uint64) (*model.Shop, error) {
	shop, err := l.svcCtx.Repo.ApproveApplication(ctx, appID, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, errors.New("申请状态不可审核")
		}
		return nil, err
	}
	return shop, nil
}

func (l *MerchantLogic) Reject(ctx context.Context, appID, adminID uint64, reason string) error {
	return l.svcCtx.Repo.RejectApplication(ctx, appID, adminID, reason)
}

func (l *MerchantLogic) ListShops(ctx context.Context, status, name string, page, pageSize int) ([]model.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListShops(ctx, status, name, page, pageSize)
}

func (l *MerchantLogic) ListPublicShops(ctx context.Context, page, pageSize int) ([]model.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return l.svcCtx.Repo.ListPublicShops(ctx, page, pageSize)
}

func (l *MerchantLogic) GetShop(ctx context.Context, id uint64) (*model.Shop, error) {
	return l.svcCtx.Repo.FindShop(ctx, id)
}

func (l *MerchantLogic) GetPublicShop(ctx context.Context, id uint64) (*model.Shop, error) {
	shop, err := l.svcCtx.Repo.FindShop(ctx, id)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}
	if shop.Status != model.ShopApproved {
		return nil, errors.New("店铺不可用")
	}
	return shop, nil
}

func (l *MerchantLogic) DisableShop(ctx context.Context, id uint64, reason string) error {
	return l.svcCtx.Repo.UpdateShopStatus(ctx, id, model.ShopDisabled, reason)
}

func (l *MerchantLogic) EnableShop(ctx context.Context, id uint64) error {
	return l.svcCtx.Repo.UpdateShopStatus(ctx, id, model.ShopApproved, "")
}

func (l *MerchantLogic) CreateShop(ctx context.Context, req types.AdminCreateShopReq) (*model.Shop, error) {
	if req.Name == "" {
		return nil, errors.New("店铺名称不能为空")
	}
	if len(req.OwnerMobile) != 11 {
		return nil, errors.New("店主手机号无效")
	}
	shop := &model.Shop{
		Name:              req.Name,
		Logo:              req.Logo,
		ContactName:       req.ContactName,
		ContactPhone:      req.ContactPhone,
		Description:       req.Description,
		Category:          req.Category,
		Province:          req.Province,
		City:              req.City,
		District:          req.District,
		Address:           req.Address,
		BusinessLicenseNo: req.BusinessLicenseNo,
		LegalPerson:       req.LegalPerson,
		LicenseImage:      req.LicenseImage,
		StorefrontImage:   req.StorefrontImage,
	}
	if shop.ContactPhone == "" {
		shop.ContactPhone = req.OwnerMobile
	}
	return l.svcCtx.Repo.CreateShopWithOwner(ctx, shop, req.OwnerMobile, req.OwnerPassword, req.OwnerNickname)
}

func (l *MerchantLogic) AdminUpdateShop(ctx context.Context, id uint64, req types.AdminUpdateShopReq) error {
	existing, err := l.svcCtx.Repo.FindShop(ctx, id)
	if err != nil {
		return errors.New("店铺不存在")
	}
	existing.Name = req.Name
	existing.Logo = req.Logo
	existing.ContactName = req.ContactName
	existing.ContactPhone = req.ContactPhone
	existing.Description = req.Description
	existing.Category = req.Category
	existing.Province = req.Province
	existing.City = req.City
	existing.District = req.District
	existing.Address = req.Address
	existing.BusinessLicenseNo = req.BusinessLicenseNo
	existing.LegalPerson = req.LegalPerson
	existing.LicenseImage = req.LicenseImage
	existing.StorefrontImage = req.StorefrontImage
	return l.svcCtx.Repo.UpdateShop(ctx, existing)
}

func (l *MerchantLogic) ResetOwnerPassword(ctx context.Context, shopID uint64, plain string) error {
	if plain == "" {
		return errors.New("密码不能为空")
	}
	return l.svcCtx.Repo.ResetOwnerPassword(ctx, shopID, plain)
}

func (l *MerchantLogic) MyShops(ctx context.Context, userID uint64) ([]model.Shop, error) {
	return l.svcCtx.Repo.ListShopsByUser(ctx, userID)
}

func (l *MerchantLogic) UpdateMyShop(ctx context.Context, shopID, userID uint64, req types.UpdateShopReq) error {
	if !l.svcCtx.Repo.IsShopMember(ctx, shopID, userID) {
		return errors.New("无权操作该店铺")
	}
	existing, err := l.svcCtx.Repo.FindShop(ctx, shopID)
	if err != nil {
		return errors.New("店铺不存在")
	}
	existing.Name = req.Name
	existing.Logo = req.Logo
	existing.ContactName = req.ContactName
	existing.ContactPhone = req.ContactPhone
	existing.Description = req.Description
	existing.Category = req.Category
	existing.Province = req.Province
	existing.City = req.City
	existing.District = req.District
	existing.Address = req.Address
	existing.StorefrontImage = req.StorefrontImage
	return l.svcCtx.Repo.UpdateShopDisplay(ctx, existing)
}
