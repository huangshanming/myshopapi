package config

import (
	"context"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConfigsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListConfigsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConfigsLogic {
	return &ListConfigsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListConfigsLogic) ListConfigs(ctx context.Context) (resp *types.PageListResp, err error) {
	data, err := hadmin.NewAdminHandler(l.svcCtx).ListConfigs(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.PageListResp{List: data}, nil
}
