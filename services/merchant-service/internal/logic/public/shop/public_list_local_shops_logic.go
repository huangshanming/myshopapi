package shop

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicListLocalShopsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicListLocalShopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListLocalShopsLogic {
	return &PublicListLocalShopsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicListLocalShopsLogic) PublicListLocalShops(ctx context.Context, req *types.LocalShopsReq) (resp *types.PageListResp, err error) {
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListLocalShops(ctx, req.Lat, req.Lng, req.Keyword, req.Sort, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
