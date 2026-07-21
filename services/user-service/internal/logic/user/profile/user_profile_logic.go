package profile

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserProfileLogic) UserProfile(ctx context.Context) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/profile", nil, nil, nil, huser.NewUserHandler(l.svcCtx).Profile)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
