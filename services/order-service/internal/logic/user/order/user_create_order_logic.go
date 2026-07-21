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

type UserCreateOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateOrderLogic {
	return &UserCreateOrderLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserCreateOrderLogic) UserCreateOrder(ctx context.Context, req *types.CreateOrderReq) (*types.AnyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	order, err := biz.NewOrderLogic(l.svcCtx).CreateOrder(ctx, userID, req.AddressID, req.Items, req.UserCouponID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: order}, nil
}
