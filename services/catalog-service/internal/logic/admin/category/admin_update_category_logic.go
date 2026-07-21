package category

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	ptypes "mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateCategoryLogic {
	return &AdminUpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateCategoryLogic) AdminUpdateCategory(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "分类ID无效")
	}
	var body ptypes.CategoryReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if body.Name == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := plogic.NewCatalogLogic(l.svcCtx).UpdateCategory(ctx, id, body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
