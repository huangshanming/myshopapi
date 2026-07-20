package article

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/public/article"
	"mymall/services/catalog-service/internal/svc"
)

func ListCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewListCommentsLogic(r.Context(), svcCtx)
		l.ListComments(w, r)
	}
}

func ListEmojisHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewListEmojisLogic(r.Context(), svcCtx)
		l.ListEmojis(w, r)
	}
}

func PublicGetArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewPublicGetArticleLogic(r.Context(), svcCtx)
		l.PublicGetArticle(w, r)
	}
}

func PublicListArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewPublicListArticlesLogic(r.Context(), svcCtx)
		l.PublicListArticles(w, r)
	}
}
