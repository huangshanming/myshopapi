package seckill

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListSeckillEntriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListSeckillEntriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListSeckillEntriesLogic {
	return &MerchantListSeckillEntriesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListSeckillEntriesLogic) MerchantListSeckillEntries(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	p, ps := req.Page, req.PageSize
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListShopSeckillEntries(shopID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
