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

type AdminCreateArticleCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateArticleCategoryLogic {
	return &AdminCreateArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateArticleCategoryLogic) AdminCreateArticleCategory(ctx context.Context, req *types.ArticleCategorySaveReq) (resp *types.EmptyResp, err error) {
	if err := clogic.NewArticleLogic(l.svcCtx).SaveCategory(ctx, 0, req.ToContent()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
