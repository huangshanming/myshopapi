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

type AdminGetShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetShopLogic {
	return &AdminGetShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetShopLogic) AdminGetShop(ctx context.Context, req *types.IdPathReq) (resp *types.ShopResp, err error) {
	id := req.Id
	shop, err := biz.NewMerchantLogic(l.svcCtx).GetShop(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "店铺不存在")
	}
	return &types.ShopResp{Data: shop}, nil
}
