package logic

import (
	"context"
	"errors"
	"strings"

	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressLogic {
	return &AddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressLogic) validate(req types.AddressReq) error {
	if strings.TrimSpace(req.ReceiverName) == "" {
		return errors.New("请填写收货人")
	}
	if strings.TrimSpace(req.ReceiverPhone) == "" {
		return errors.New("请填写手机号")
	}
	if strings.TrimSpace(req.Province) == "" || strings.TrimSpace(req.City) == "" || strings.TrimSpace(req.District) == "" {
		return errors.New("请选择省市区")
	}
	if strings.TrimSpace(req.ProvinceCode) == "" || strings.TrimSpace(req.CityCode) == "" || strings.TrimSpace(req.DistrictCode) == "" {
		return errors.New("请选择省市区")
	}
	if strings.TrimSpace(req.Detail) == "" {
		return errors.New("请填写详细地址")
	}
	return nil
}

func (l *AddressLogic) List(userID uint64) ([]model.UserAddress, error) {
	if userID == 0 {
		return nil, errors.New("用户无效")
	}
	return l.svcCtx.Repo.ListAddresses(userID)
}

func (l *AddressLogic) Get(userID, id uint64) (*model.UserAddress, error) {
	if userID == 0 || id == 0 {
		return nil, errors.New("参数无效")
	}
	return l.svcCtx.Repo.GetAddressByID(userID, id)
}

func (l *AddressLogic) Create(userID uint64, req types.AddressReq) (*model.UserAddress, error) {
	if userID == 0 {
		return nil, errors.New("用户无效")
	}
	if err := l.validate(req); err != nil {
		return nil, err
	}
	a := &model.UserAddress{
		UserID:        userID,
		ReceiverName:  strings.TrimSpace(req.ReceiverName),
		ReceiverPhone: strings.TrimSpace(req.ReceiverPhone),
		Province:      strings.TrimSpace(req.Province),
		City:          strings.TrimSpace(req.City),
		District:      strings.TrimSpace(req.District),
		Detail:        strings.TrimSpace(req.Detail),
		ProvinceCode:  strings.TrimSpace(req.ProvinceCode),
		CityCode:      strings.TrimSpace(req.CityCode),
		DistrictCode:  strings.TrimSpace(req.DistrictCode),
		IsDefault:     req.IsDefault,
	}
	if a.IsDefault != 1 {
		a.IsDefault = 0
	}
	if err := l.svcCtx.Repo.CreateAddress(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (l *AddressLogic) Update(userID, id uint64, req types.AddressReq) error {
	if err := l.validate(req); err != nil {
		return err
	}
	a := &model.UserAddress{
		ReceiverName:  strings.TrimSpace(req.ReceiverName),
		ReceiverPhone: strings.TrimSpace(req.ReceiverPhone),
		Province:      strings.TrimSpace(req.Province),
		City:          strings.TrimSpace(req.City),
		District:      strings.TrimSpace(req.District),
		Detail:        strings.TrimSpace(req.Detail),
		ProvinceCode:  strings.TrimSpace(req.ProvinceCode),
		CityCode:      strings.TrimSpace(req.CityCode),
		DistrictCode:  strings.TrimSpace(req.DistrictCode),
		IsDefault:     req.IsDefault,
	}
	if a.IsDefault != 1 {
		a.IsDefault = 0
	}
	return l.svcCtx.Repo.UpdateAddress(userID, id, a)
}

func (l *AddressLogic) Delete(userID, id uint64) error {
	return l.svcCtx.Repo.DeleteAddress(userID, id)
}

func (l *AddressLogic) SetDefault(userID, id uint64) error {
	return l.svcCtx.Repo.SetDefaultAddress(userID, id)
}
