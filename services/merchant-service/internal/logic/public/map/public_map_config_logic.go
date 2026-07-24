package mapapi

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicMapConfigLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicMapConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicMapConfigLogic {
	return &PublicMapConfigLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicMapConfigLogic) PublicMapConfig(ctx context.Context) (resp *types.AnyResp, err error) {
	key := ""
	if l.svcCtx.Config != nil {
		key = l.svcCtx.Config.TencentMap.Key
	}
	if key == "" {
		return nil, xerr.New(http.StatusBadRequest, "未配置腾讯地图 Key")
	}
	return &types.AnyResp{Data: map[string]string{"key": key}}, nil
}
