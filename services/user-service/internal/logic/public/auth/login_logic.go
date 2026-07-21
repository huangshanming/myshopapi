package auth

import (
	"context"
	"net/http"
	"github.com/zeromicro/go-zero/core/logx"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type LoginLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *LoginLogic) Login(ctx context.Context, req *types.LoginReq) (*types.LoginResp, error) {
	token, user, err := biz.NewUserLogic(l.svcCtx).LoginWithShop(ctx, req.Mobile, req.Password, req.ShopID)
	if err != nil {
		return nil, xerr.New(http.StatusUnauthorized, err.Error())
	}
	return &types.LoginResp{
		Token: token,
		User: map[string]interface{}{
			"id": user.ID, "mobile": user.Mobile, "nickname": user.Nickname,
			"avatar": user.Avatar, "role": user.Role, "status": user.Status,
		},
	}, nil
}
