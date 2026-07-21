package profile

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserProfileLogic) UserProfile(ctx context.Context) (resp *types.AnyResp, err error) {
	data, err := huser.NewUserHandler(l.svcCtx).Profile(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
