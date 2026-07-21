package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelScheduleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCancelScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelScheduleLogic {
	return &CancelScheduleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CancelScheduleLogic) CancelSchedule(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id := req.Id
	_ = l.svcCtx.ProductAdmin.CancelSchedule(ctx, id, shopID)
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
