package logic

import (
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

func (l *MerchantLogic) Apply(userID uint64, in types.ApplyReq) (*model.ShopApplication, error) {
	if _, err := l.svcCtx.Repo.FindPendingAppByUser(userID); err == nil {
		return nil, errors.New("已有待审核的入驻申请")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	app := &model.ShopApplication{
		UserID:       userID,
		ShopName:     in.ShopName,
		ContactName:  in.ContactName,
		ContactPhone: in.ContactPhone,
		Description:  in.Description,
		Status:       model.AppPending,
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

func (l *MerchantLogic) ListShops(status string, page, pageSize int) ([]model.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return l.svcCtx.Repo.ListShops(status, page, pageSize)
}

func (l *MerchantLogic) GetShop(id uint64) (*model.Shop, error) {
	return l.svcCtx.Repo.FindShop(id)
}

func (l *MerchantLogic) DisableShop(id uint64, reason string) error {
	return l.svcCtx.Repo.UpdateShopStatus(id, model.ShopDisabled, reason)
}

func (l *MerchantLogic) EnableShop(id uint64) error {
	return l.svcCtx.Repo.UpdateShopStatus(id, model.ShopApproved, "")
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
	return l.svcCtx.Repo.UpdateShop(existing)
}
