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

type PublicGetShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicGetShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicGetShopLogic {
	return &PublicGetShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicGetShopLogic) PublicGetShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	shop, err := biz.NewMerchantLogic(l.svcCtx).GetPublicShop(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: shop}, nil
}
