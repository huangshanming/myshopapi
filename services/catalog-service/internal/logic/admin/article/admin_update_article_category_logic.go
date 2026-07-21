package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateArticleCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateArticleCategoryLogic {
	return &AdminUpdateArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateArticleCategoryLogic) AdminUpdateArticleCategory(ctx context.Context, req *types.ArticleCategoryUpdateBodyReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).SaveCategory(ctx, id, req.ToContent()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
