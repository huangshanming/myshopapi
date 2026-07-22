package profile

import (
	"context"
	"net/http"
	"github.com/zeromicro/go-zero/core/logx"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type UserProfileLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserProfileLogic) UserProfile(ctx context.Context) (*types.UserResp, error) {
	var userID uint64
	if id, ok := middleware.GetUserID(ctx); ok {
		userID = id
	} else if claims, ok := jwt.ClaimsFromContext(ctx); ok {
		userID = claims.UserID
	} else {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	user, err := biz.NewUserLogic(l.svcCtx).GetProfile(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "用户不存在")
	}
	return &types.UserResp{Data: user}, nil
}
