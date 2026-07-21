package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"net/http"

	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUpdateTagLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUpdateTagLogic {
	return &MerchantUpdateTagLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUpdateTagLogic) MerchantUpdateTag(ctx context.Context, req *types.TagUpdateBodyReq) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	tag := &model.ProductTag{ID: req.Id, ShopID: shopID, Name: req.Name, Color: req.Color, Status: 1}
	if err := l.svcCtx.ProductAdmin.SaveTag(ctx, tag); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: tag}, nil
}
