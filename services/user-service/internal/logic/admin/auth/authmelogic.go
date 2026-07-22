// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthMeLogic {
	return &AuthMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthMeLogic) AuthMe() (resp *types.AuthMeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
