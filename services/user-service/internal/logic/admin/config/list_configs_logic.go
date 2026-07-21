package config

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
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
	list, err := biz.NewRBACLogic(l.svcCtx).ListConfigs(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
