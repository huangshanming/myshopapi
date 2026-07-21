package config

import (
	"context"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *SaveConfigsLogic) SaveConfigs(ctx context.Context, req *types.ConfigBatchReq) error {
	_, err := hadmin.NewAdminHandler(l.svcCtx).SaveConfigs(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return err
	}
	return nil
}
