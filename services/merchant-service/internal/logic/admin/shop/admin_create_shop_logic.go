package shop

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateShopLogic {
	return &AdminCreateShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateShopLogic) AdminCreateShop(ctx context.Context, req *types.AdminCreateShopReq) (resp *types.ShopResp, err error) {
	shop, err := biz.NewMerchantLogic(l.svcCtx).CreateShop(ctx, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.ShopResp{Data: shop}, nil
}
