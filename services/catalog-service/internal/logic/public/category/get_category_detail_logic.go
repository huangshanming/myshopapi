package category

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

type GetCategoryDetailLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetCategoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryDetailLogic {
	return &GetCategoryDetailLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryDetailLogic) GetCategoryDetail(ctx context.Context) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{}

	id, _ := strconv.ParseUint(in.QueryGet("id"), 10, 64)
	if id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetCategoryDetail(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerr.New(http.StatusNotFound, "分类不存在")
		}
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.AnyResp{Data: data}, nil
}
