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

type RegisterLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(ctx context.Context, req *types.RegisterReq) (*types.UserResp, error) {
	user, err := biz.NewUserLogic(l.svcCtx).Register(ctx, req.Mobile, req.Password)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.UserResp{Data: user}, nil
}
