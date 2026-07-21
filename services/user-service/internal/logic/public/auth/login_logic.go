package auth

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(ctx context.Context, req *types.LoginReq) (resp *types.LoginResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/login", nil, nil, req, huser.NewUserHandler(l.svcCtx).Login)
	if err != nil {
		return nil, err
	}
	var out types.LoginResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
