package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ListUserPointsOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListUserPointsOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserPointsOrdersLogic {
	return &ListUserPointsOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListUserPointsOrdersLogic) ListUserPointsOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewPointsOrderLogic(l.svcCtx).UserList(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil
}
