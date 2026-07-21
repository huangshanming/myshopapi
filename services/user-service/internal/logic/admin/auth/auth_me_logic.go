package auth

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAuthMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthMeLogic {
	return &AuthMeLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AuthMeLogic) AuthMe(ctx context.Context) (*types.AuthMeResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	data, err := biz.NewRBACLogic(l.svcCtx).AuthMe(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}
