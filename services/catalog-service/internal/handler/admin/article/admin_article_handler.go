package article

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/article"
	"mymall/services/catalog-service/internal/svc"
)

func AdminArticleStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminArticleStatsLogic(r.Context(), svcCtx)
		l.AdminArticleStats(w, r)
	}
}

func AdminAuditArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminAuditArticleLogic(r.Context(), svcCtx)
		l.AdminAuditArticle(w, r)
	}
}

func AdminBatchAuditArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminBatchAuditArticlesLogic(r.Context(), svcCtx)
		l.AdminBatchAuditArticles(w, r)
	}
}

func AdminCreateArticleCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminCreateArticleCategoryLogic(r.Context(), svcCtx)
		l.AdminCreateArticleCategory(w, r)
	}
}

func AdminCreateArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminCreateArticleLogic(r.Context(), svcCtx)
		l.AdminCreateArticle(w, r)
	}
}

func AdminDeleteArticleCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminDeleteArticleCategoryLogic(r.Context(), svcCtx)
		l.AdminDeleteArticleCategory(w, r)
	}
}

func AdminGetArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminGetArticleLogic(r.Context(), svcCtx)
		l.AdminGetArticle(w, r)
	}
}

func AdminListArticleCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminListArticleCategoriesLogic(r.Context(), svcCtx)
		l.AdminListArticleCategories(w, r)
	}
}

func AdminListArticleRecycleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminListArticleRecycleLogic(r.Context(), svcCtx)
		l.AdminListArticleRecycle(w, r)
	}
}

func AdminListArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminListArticlesLogic(r.Context(), svcCtx)
		l.AdminListArticles(w, r)
	}
}

func AdminOfflineArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminOfflineArticleLogic(r.Context(), svcCtx)
		l.AdminOfflineArticle(w, r)
	}
}

func AdminPurgeArticleRecycleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminPurgeArticleRecycleLogic(r.Context(), svcCtx)
		l.AdminPurgeArticleRecycle(w, r)
	}
}

func AdminRestoreArticleRecycleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminRestoreArticleRecycleLogic(r.Context(), svcCtx)
		l.AdminRestoreArticleRecycle(w, r)
	}
}

func AdminSoftDeleteArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminSoftDeleteArticleLogic(r.Context(), svcCtx)
		l.AdminSoftDeleteArticle(w, r)
	}
}

func AdminTopArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminTopArticleLogic(r.Context(), svcCtx)
		l.AdminTopArticle(w, r)
	}
}

func AdminUpdateArticleCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminUpdateArticleCategoryLogic(r.Context(), svcCtx)
		l.AdminUpdateArticleCategory(w, r)
	}
}

func AdminUpdateArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminUpdateArticleLogic(r.Context(), svcCtx)
		l.AdminUpdateArticle(w, r)
	}
}

func AdminUploadArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewAdminUploadArticleLogic(r.Context(), svcCtx)
		l.AdminUploadArticle(w, r)
	}
}
