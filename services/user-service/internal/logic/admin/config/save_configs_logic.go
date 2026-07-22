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

type SaveConfigsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSaveConfigsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveConfigsLogic {
	return &SaveConfigsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SaveConfigsLogic) SaveConfigs(ctx context.Context, req *types.ConfigBatchReq) (*types.EmptyResp, error) {
	if err := biz.NewRBACLogic(l.svcCtx).SaveConfigs(ctx, req.Items); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
