package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StockWarningsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewStockWarningsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockWarningsLogic {
	return &StockWarningsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *StockWarningsLogic) StockWarnings(ctx context.Context, req *types.PageReq) (resp *types.StockWarningsResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	data, err := plogic.NewProductAdminLogic(l.svcCtx).StockWarnings(ctx, shopID, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.StockWarningsResp{Data: data}, nil
}
