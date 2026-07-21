package product

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *StockWarningsLogic) StockWarnings(ctx context.Context) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	data, err := plogic.NewProductAdminLogic(l.svcCtx).StockWarnings(ctx, shopID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
