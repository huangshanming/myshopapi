package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ListPointsOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListPointsOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointsOrdersLogic {
	return &ListPointsOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListPointsOrdersLogic) ListPointsOrders(ctx context.Context, req *types.ListPointsOrdersReq) (resp *types.PageListResp, err error) {
	list, total, err := biz.NewPointsOrderLogic(l.svcCtx).AdminList(ctx, req.Page, req.PageSize, req.Status, req.OrderNo, req.Keyword, req.UserID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil
}
