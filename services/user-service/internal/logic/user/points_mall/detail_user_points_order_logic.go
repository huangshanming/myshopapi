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

type DetailUserPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDetailUserPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailUserPointsOrderLogic {
	return &DetailUserPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DetailUserPointsOrderLogic) DetailUserPointsOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	o, err := biz.NewPointsOrderLogic(l.svcCtx).UserGet(ctx, userID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
