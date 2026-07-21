package config

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveConfigsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSaveConfigsLogic(svcCtx *svc.ServiceContext) *SaveConfigsLogic {
	return &SaveConfigsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *SaveConfigsLogic) SaveConfigs(ctx context.Context, req *types.ConfigBatchReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/configs", nil, nil, req, hadmin.NewAdminHandler(l.svcCtx).SaveConfigs)
	if err != nil {
		return err
	}
	return nil
}
