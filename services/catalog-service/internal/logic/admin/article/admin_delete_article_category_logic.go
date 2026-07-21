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

type AdminDeleteArticleCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteArticleCategoryLogic {
	return &AdminDeleteArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteArticleCategoryLogic) AdminDeleteArticleCategory(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).DeleteCategory(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
