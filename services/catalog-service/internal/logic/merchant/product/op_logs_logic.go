package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewOpLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpLogsLogic {
	return &OpLogsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *OpLogsLogic) OpLogs(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	pid, _ := strconv.ParseUint("" /* was query:product_id */, 10, 64)
	page, pageSize := req.Page, req.PageSize
	data, err := plogic.NewProductAdminLogic(l.svcCtx).OpLogs(ctx, shopID, pid, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
