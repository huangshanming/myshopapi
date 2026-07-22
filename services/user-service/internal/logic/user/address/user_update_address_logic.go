package address

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateAddressLogic {
	return &UserUpdateAddressLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserUpdateAddressLogic) UserUpdateAddress(ctx context.Context, req *types.AddressUpdateReq) (*types.EmptyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	if err := biz.NewAddressLogic(l.svcCtx).Update(ctx, userID, req.Id, types.AddressReq{
		ReceiverName: req.ReceiverName, ReceiverPhone: req.ReceiverPhone,
		Province: req.Province, City: req.City, District: req.District, Detail: req.Detail,
		ProvinceCode: req.ProvinceCode, CityCode: req.CityCode, DistrictCode: req.DistrictCode,
		IsDefault: req.IsDefault,
	}); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
