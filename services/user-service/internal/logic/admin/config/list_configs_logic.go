package config

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConfigsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListConfigsLogic(svcCtx *svc.ServiceContext) *ListConfigsLogic {
	return &ListConfigsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ListConfigsLogic) ListConfigs(ctx context.Context) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/configs", nil, nil, nil, hadmin.NewAdminHandler(l.svcCtx).ListConfigs)
	if err != nil {
		return nil, err
	}
	var list interface{}
	if err := httpinvoke.Decode(raw, &list); err != nil {
		return nil, err
	}
	return &types.PageListResp{List: list}, nil
}
