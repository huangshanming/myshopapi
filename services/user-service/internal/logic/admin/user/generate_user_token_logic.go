package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type GenerateUserTokenLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGenerateUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateUserTokenLogic {
	return &GenerateUserTokenLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GenerateUserTokenLogic) GenerateUserToken(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	data, err := biz.NewRBACLogic(l.svcCtx).GenerateUserToken(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
