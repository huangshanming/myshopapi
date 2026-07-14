package service

import (
	"errors"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/repository"

	"gorm.io/gorm"
)

type MerchantService struct {
	repo *repository.MerchantRepository
}

func NewMerchantService(repo *repository.MerchantRepository) *MerchantService {
	return &MerchantService{repo: repo}
}

type ApplyInput struct {
	ShopName     string `json:"shop_name" binding:"required"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required,len=11"`
	Description  string `json:"description"`
}

func (s *MerchantService) Apply(userID uint64, in ApplyInput) (*model.ShopApplication, error) {
	if _, err := s.repo.FindPendingAppByUser(userID); err == nil {
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
	if err := s.repo.CreateApplication(app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *MerchantService) ListApplications(status string, page, pageSize int) ([]model.ShopApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListApplications(status, page, pageSize)
}

func (s *MerchantService) Approve(appID, adminID uint64) (*model.Shop, error) {
	shop, err := s.repo.ApproveApplication(appID, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, errors.New("申请状态不可审核")
		}
		return nil, err
	}
	return shop, nil
}

func (s *MerchantService) Reject(appID, adminID uint64, reason string) error {
	return s.repo.RejectApplication(appID, adminID, reason)
}

func (s *MerchantService) ListShops(status string, page, pageSize int) ([]model.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListShops(status, page, pageSize)
}

func (s *MerchantService) GetShop(id uint64) (*model.Shop, error) {
	return s.repo.FindShop(id)
}

func (s *MerchantService) DisableShop(id uint64, reason string) error {
	return s.repo.UpdateShopStatus(id, model.ShopDisabled, reason)
}

func (s *MerchantService) EnableShop(id uint64) error {
	return s.repo.UpdateShopStatus(id, model.ShopApproved, "")
}

func (s *MerchantService) MyShops(userID uint64) ([]model.Shop, error) {
	return s.repo.ListShopsByUser(userID)
}

func (s *MerchantService) UpdateMyShop(shopID, userID uint64, shop *model.Shop) error {
	if !s.repo.IsShopMember(shopID, userID) {
		return errors.New("无权操作该店铺")
	}
	existing, err := s.repo.FindShop(shopID)
	if err != nil {
		return errors.New("店铺不存在")
	}
	existing.Name = shop.Name
	existing.Logo = shop.Logo
	existing.ContactName = shop.ContactName
	existing.ContactPhone = shop.ContactPhone
	existing.Description = shop.Description
	return s.repo.UpdateShop(existing)
}
