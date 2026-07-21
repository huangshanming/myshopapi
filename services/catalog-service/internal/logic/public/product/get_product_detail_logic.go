package product

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

func (l *GetProductDetailLogic) GetProductDetail(ctx context.Context) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{}

	id, _ := strconv.ParseUint(in.QueryGet("id"), 10, 64)
	if id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetProductDetail(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerr.New(http.StatusNotFound, "商品不存在")
		}
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.AnyResp{Data: data}, nil
}
