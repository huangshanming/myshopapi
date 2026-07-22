package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicGetArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicGetArticleLogic {
	return &PublicGetArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicGetArticleLogic) PublicGetArticle(ctx context.Context, req *types.IdPathReq) (resp *types.ArticleResp, err error) {
	userID, _ := middleware.GetUserID(ctx)
	data, err := clogic.NewArticleLogic(l.svcCtx).PublicDetail(ctx, req.Id, userID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.ArticleResp{Data: data}, nil
}
