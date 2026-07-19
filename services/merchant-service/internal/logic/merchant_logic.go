package logic

import (
	"context"
	"errors"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"gorm.io/gorm"
)

type MerchantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantLogic {
	return &MerchantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantLogic) Apply(userID uint64, in types.ApplyReq) (*model.ShopApplication, error) {
	if _, err := l.svcCtx.Repo.FindPendingAppByUser(userID); err == nil {
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
	if err := l.svcCtx.Repo.CreateApplication(app); err != nil {
		return nil, err
	}
	return app, nil
}

func (l *MerchantLogic) ListApplications(status string, page, pageSize int) ([]model.ShopApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListApplications(status, page, pageSize)
}

func (l *MerchantLogic) Approve(appID, adminID uint64) (*model.Shop, error) {
	shop, err := l.svcCtx.Repo.ApproveApplication(appID, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, errors.New("申请状态不可审核")
		}
		return nil, err
	}
	return shop, nil
}

func (l *MerchantLogic) Reject(appID, adminID uint64, reason string) error {
	return l.svcCtx.Repo.RejectApplication(appID, adminID, reason)
}

func (l *MerchantLogic) ListShops(status, name string, page, pageSize int) ([]model.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListShops(status, name, page, pageSize)
}

func (l *MerchantLogic) ListPublicShops(page, pageSize int) ([]model.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return l.svcCtx.Repo.ListPublicShops(page, pageSize)
}

func (l *MerchantLogic) GetShop(id uint64) (*model.Shop, error) {
	return l.svcCtx.Repo.FindShop(id)
}

func (l *MerchantLogic) GetPublicShop(id uint64) (*model.Shop, error) {
	shop, err := l.svcCtx.Repo.FindShop(id)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}
	if shop.Status != model.ShopApproved {
		return nil, errors.New("店铺不可用")
	}
	return shop, nil
}

func (l *MerchantLogic) DisableShop(id uint64, reason string) error {
	return l.svcCtx.Repo.UpdateShopStatus(id, model.ShopDisabled, reason)
}

func (l *MerchantLogic) EnableShop(id uint64) error {
	return l.svcCtx.Repo.UpdateShopStatus(id, model.ShopApproved, "")
}

func (l *MerchantLogic) CreateShop(req types.AdminCreateShopReq) (*model.Shop, error) {
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
	return l.svcCtx.Repo.CreateShopWithOwner(shop, req.OwnerMobile, req.OwnerPassword, req.OwnerNickname)
}

func (l *MerchantLogic) AdminUpdateShop(id uint64, req types.AdminUpdateShopReq) error {
	existing, err := l.svcCtx.Repo.FindShop(id)
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
	return l.svcCtx.Repo.UpdateShop(existing)
}

func (l *MerchantLogic) ResetOwnerPassword(shopID uint64, plain string) error {
	if plain == "" {
		return errors.New("密码不能为空")
	}
	return l.svcCtx.Repo.ResetOwnerPassword(shopID, plain)
}

func (l *MerchantLogic) MyShops(userID uint64) ([]model.Shop, error) {
	return l.svcCtx.Repo.ListShopsByUser(userID)
}

func (l *MerchantLogic) UpdateMyShop(shopID, userID uint64, req types.UpdateShopReq) error {
	if !l.svcCtx.Repo.IsShopMember(shopID, userID) {
		return errors.New("无权操作该店铺")
	}
	existing, err := l.svcCtx.Repo.FindShop(shopID)
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
	return l.svcCtx.Repo.UpdateShopDisplay(existing)
}
