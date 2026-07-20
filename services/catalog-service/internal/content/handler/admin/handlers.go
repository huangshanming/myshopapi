package admin

import (
	"context"

	"mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/svc"
)

type ArticleHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.ArticleLogic
}

func NewArticleHandler(svcCtx *svc.ServiceContext) *ArticleHandler {
	return &ArticleHandler{svcCtx: svcCtx, logic: logic.NewArticleLogic(context.Background(), svcCtx)}
}
