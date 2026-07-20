package article

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/merchant/article"
	"mymall/services/catalog-service/internal/svc"
)

func MerchantCreateArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantCreateArticleLogic(r.Context(), svcCtx)
		l.MerchantCreateArticle(w, r)
	}
}

func MerchantDeleteArticleCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantDeleteArticleCommentLogic(r.Context(), svcCtx)
		l.MerchantDeleteArticleComment(w, r)
	}
}

func MerchantDeleteArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantDeleteArticleLogic(r.Context(), svcCtx)
		l.MerchantDeleteArticle(w, r)
	}
}

func MerchantGetArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantGetArticleLogic(r.Context(), svcCtx)
		l.MerchantGetArticle(w, r)
	}
}

func MerchantListArticleCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantListArticleCategoriesLogic(r.Context(), svcCtx)
		l.MerchantListArticleCategories(w, r)
	}
}

func MerchantListArticleCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantListArticleCommentsLogic(r.Context(), svcCtx)
		l.MerchantListArticleComments(w, r)
	}
}

func MerchantListArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantListArticlesLogic(r.Context(), svcCtx)
		l.MerchantListArticles(w, r)
	}
}

func MerchantPatchArticleCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantPatchArticleCommentLogic(r.Context(), svcCtx)
		l.MerchantPatchArticleComment(w, r)
	}
}

func MerchantUpdateArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantUpdateArticleLogic(r.Context(), svcCtx)
		l.MerchantUpdateArticle(w, r)
	}
}

func MerchantUploadArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewMerchantUploadArticleLogic(r.Context(), svcCtx)
		l.MerchantUploadArticle(w, r)
	}
}
