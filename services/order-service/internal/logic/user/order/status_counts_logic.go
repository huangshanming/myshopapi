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

type StatusCountsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewStatusCountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusCountsLogic {
	return &StatusCountsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *StatusCountsLogic) StatusCounts(ctx context.Context) (*types.OrderStatusCountsResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	counts, err := biz.NewOrderLogic(l.svcCtx).UserOrderStatusCounts(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.OrderStatusCountsResp{Counts: counts}, nil
}
