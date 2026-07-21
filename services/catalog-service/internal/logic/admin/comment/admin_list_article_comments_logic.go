package comment

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListArticleCommentsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListArticleCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticleCommentsLogic {
	return &AdminListArticleCommentsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListArticleCommentsLogic) AdminListArticleComments(ctx context.Context, req *types.ArticleCommentListReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	articleID := req.ArticleId
	shopID := req.ShopId
	data, err := clogic.NewArticleLogic(l.svcCtx).ListComments(ctx, repository.CommentListFilter{
		ShopID: shopID, ArticleID: articleID, Status: req.Status,
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
