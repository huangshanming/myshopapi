package points_mall

import (
	"context"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductLogic {
	return &ProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ProductLogic) ListProducts(ctx context.Context, req *types.ProductListReq) (*types.ProductListResp, error) {
	list, total, err := biz.NewPointsProductLogic(l.svcCtx).List(ctx, req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		return nil, err
	}
	return &types.ProductListResp{
		Total: total,
		List:  list,
	}, nil
}

func (l *ProductLogic) DetailProduct(ctx context.Context, req *types.IdPathReq) (*types.ProductDetailResp, error) {
	product, err := biz.NewPointsProductLogic(l.svcCtx).Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &types.ProductDetailResp{
		Product: product,
	}, nil
}
