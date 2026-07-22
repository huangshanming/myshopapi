package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ExchangeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ExchangeLogic) Exchange(ctx context.Context, req *types.ExchangeReq) (resp *types.PointsOrderResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	bizReq := biz.ExchangeReq{
		ProductID: req.ProductID,
		Quantity:  1,
	}
	if req.AddressID > 0 {
		addr, err := biz.NewAddressLogic(l.svcCtx).Get(ctx, userID, req.AddressID)
		if err != nil {
			return nil, xerr.New(http.StatusBadRequest, err.Error())
		}
		bizReq.ReceiverName = addr.ReceiverName
		bizReq.ReceiverPhone = addr.ReceiverPhone
		bizReq.ReceiverAddress = addr.FullAddress()
	}
	o, err := biz.NewPointsOrderLogic(l.svcCtx).UserExchange(ctx, userID, bizReq)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsOrderResp{Data: o}, nil
}
