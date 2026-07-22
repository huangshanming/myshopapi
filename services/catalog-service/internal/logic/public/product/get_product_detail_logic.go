package product

import (
	"context"
	"errors"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetProductDetailLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetProductDetailLogic) GetProductDetail(ctx context.Context, req *types.IdQueryReq) (resp *types.ProductResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetProductDetail(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, xerr.New(http.StatusNotFound, "商品不存在")
		}
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.ProductResp{Data: data}, nil
}
