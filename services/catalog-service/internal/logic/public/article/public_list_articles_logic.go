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

type PublicListArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListArticlesLogic {
	return &PublicListArticlesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicListArticlesLogic) PublicListArticles(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := int(req.Page), int(req.PageSize)
	home := "" /* was query:home */ == "1"
	data, err := clogic.NewArticleLogic(l.svcCtx).PublicList(ctx, page, pageSize, home)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
