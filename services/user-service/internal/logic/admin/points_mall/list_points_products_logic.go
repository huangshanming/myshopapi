package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ListPointsProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListPointsProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointsProductsLogic {
	return &ListPointsProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListPointsProductsLogic) ListPointsProducts(ctx context.Context, req *types.ListPointsProductsReq) (resp *types.PageListResp, err error) {
	list, total, err := biz.NewPointsProductLogic(l.svcCtx).List(ctx, req.Page, req.PageSize, req.Status, req.Keyword)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil
}
