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

type BatchStockLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewBatchStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchStockLogic {
	return &BatchStockLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *BatchStockLogic) BatchStock(ctx context.Context, req *types.BatchStockReq) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if err := plogic.NewProductAdminLogic(l.svcCtx).BatchStock(ctx, shopID, req.ToProduct()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
