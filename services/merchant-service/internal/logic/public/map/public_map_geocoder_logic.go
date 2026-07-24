package mapapi

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicMapGeocoderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicMapGeocoderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicMapGeocoderLogic {
	return &PublicMapGeocoderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicMapGeocoderLogic) PublicMapGeocoder(ctx context.Context, req *types.MapGeocoderReq) (resp *types.AnyResp, err error) {
	if req.Lat == 0 && req.Lng == 0 {
		return nil, xerr.New(http.StatusBadRequest, "缺少坐标")
	}
	if l.svcCtx.TencentMap == nil || !l.svcCtx.TencentMap.Configured() {
		return nil, xerr.New(http.StatusBadRequest, "未配置腾讯地图 Key")
	}
	res, err := l.svcCtx.TencentMap.ReverseGeocode(req.Lat, req.Lng)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: res}, nil
}
