package product

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountLogic {
	return &CountLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CountLogic) Count(ctx context.Context, req *types.IdPathReq) (resp *types.CountResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	n, err := plogic.NewFavoriteLogic(l.svcCtx).FavoriteCount(req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "商品不存在")
	}
	return &types.CountResp{Count: n}, nil
}
