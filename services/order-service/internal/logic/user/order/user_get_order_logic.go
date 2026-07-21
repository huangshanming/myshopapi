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

type UserGetOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetOrderLogic {
	return &UserGetOrderLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserGetOrderLogic) UserGetOrder(ctx context.Context, req *types.IdPathReq) (*types.AnyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	order, err := biz.NewOrderLogic(l.svcCtx).GetOrder(ctx, userID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "订单不存在")
	}
	return &types.AnyResp{Data: order}, nil
}
