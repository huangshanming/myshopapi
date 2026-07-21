package comment

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"
	"strconv"

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

func (l *AdminListArticleCommentsLogic) AdminListArticleComments(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	articleID, _ := strconv.ParseUint("" /* was query:article_id */, 10, 64)
	shopID, _ := strconv.ParseUint("" /* was query:shop_id */, 10, 64)
	data, err := clogic.NewArticleLogic(l.svcCtx).ListComments(ctx, repository.CommentListFilter{
		ShopID: shopID, ArticleID: articleID, Status: "" /* was query:status */,
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
