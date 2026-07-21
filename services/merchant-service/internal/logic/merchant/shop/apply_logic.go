package shop

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ApplyLogic) Apply(ctx context.Context, req *types.ApplyReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	app, err := biz.NewMerchantLogic(l.svcCtx).Apply(ctx, userID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: app}, nil
}
