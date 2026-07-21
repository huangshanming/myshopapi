package auth

import (
	"context"
	"encoding/json"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(ctx context.Context, req *types.LoginReq) (resp *types.LoginResp, err error) {
	data, err := huser.NewUserHandler(l.svcCtx).Login(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.LoginResp
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
