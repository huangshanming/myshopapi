package order

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmReceiveLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewConfirmReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmReceiveLogic {
	return &ConfirmReceiveLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *ConfirmReceiveLogic) ConfirmReceive(ctx context.Context, req *types.IdPathReq) (*types.EmptyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	if err := biz.NewOrderLogic(l.svcCtx).ConfirmReceive(ctx, userID, req.Id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
