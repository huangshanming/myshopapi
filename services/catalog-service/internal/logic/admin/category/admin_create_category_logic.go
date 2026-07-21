package category

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	ptypes "mymall/services/catalog-service/internal/product/types"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateCategoryLogic {
	return &AdminCreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCategoryLogic) AdminCreateCategory(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body ptypes.CategoryReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if body.Name == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	cat, err := plogic.NewCatalogLogic(l.svcCtx).CreateCategory(ctx, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: cat}, nil
}
