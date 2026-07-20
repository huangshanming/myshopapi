package handler

import (
	"net/http"

	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/metrics"
	cadmin "mymall/services/catalog-service/internal/content/handler/admin"
	cmerchant "mymall/services/catalog-service/internal/content/handler/merchant"
	cpublic "mymall/services/catalog-service/internal/content/handler/public"
	svcMW "mymall/services/catalog-service/internal/middleware"
	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	padmin "mymall/services/catalog-service/internal/product/handler/admin"
	pmerchant "mymall/services/catalog-service/internal/product/handler/merchant"
	ppublic "mymall/services/catalog-service/internal/product/handler/public"
	puser "mymall/services/catalog-service/internal/product/handler/user"
	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/uploadpath"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext, healthReg *health.Registry, mws svcMW.Bundle) {
	catalogPublic := ppublic.NewCatalogHandler(svcCtx)
	catalogAdmin := padmin.NewCatalogHandler(svcCtx)
	favoriteUser := puser.NewFavoriteHandler(svcCtx)
	favoriteAdmin := padmin.NewFavoriteHandler(svcCtx)
	adminH := pmerchant.NewProductHandler(svcCtx)
	shopOpsH := shopopshandler.NewShopOpsHandler(svcCtx)
	articleAdminH := cadmin.NewArticleHandler(svcCtx)
	articleMerchantH := cmerchant.NewArticleHandler(svcCtx)
	articlePublicH := cpublic.NewArticleHandler(svcCtx)
	shopUploadH := padmin.NewShopUploadHandler()
	platformProductH := padmin.NewPlatformProductHandler(svcCtx)
	notifH := notifyhandler.NewNotificationHandler(svcCtx)

	server.AddRoutes(mws.Public([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: httpserver.Healthz("catalog-service")},
		{Method: http.MethodGet, Path: "/readyz", Handler: healthReg.ReadyHandler()},
		{Method: http.MethodGet, Path: "/metrics", Handler: metrics.Handler()},
		{Method: http.MethodGet, Path: "/api/v1/products/list", Handler: catalogPublic.GetProductList},
		{Method: http.MethodGet, Path: "/api/v1/products/detail", Handler: catalogPublic.GetProductDetail},
		{Method: http.MethodGet, Path: "/api/v1/products/sales-rank", Handler: catalogPublic.GetSalesRank},
		{Method: http.MethodGet, Path: "/api/v1/product_category/list", Handler: catalogPublic.GetCategoryList},
		{Method: http.MethodGet, Path: "/api/v1/product_category/detail", Handler: catalogPublic.GetCategoryDetail},
		{Method: http.MethodGet, Path: "/api/v1/banners", Handler: articlePublicH.ListBanners},
		{Method: http.MethodGet, Path: "/api/v1/articles/list", Handler: articlePublicH.List},
		{Method: http.MethodGet, Path: "/api/v1/articles/:id", Handler: articlePublicH.Detail},
		{Method: http.MethodGet, Path: "/api/v1/articles/:id/comments", Handler: articlePublicH.ListComments},
		{Method: http.MethodGet, Path: "/api/v1/comment-emojis", Handler: articlePublicH.ListEmojis},
		{Method: http.MethodGet, Path: "/api/v1/products/:id/favorite-count", Handler: favoriteUser.Count},
	}))
	server.AddRoutes(mws.Authed([]rest.Route{
		{Method: http.MethodPost, Path: "/api/v1/articles/:id/comments", Handler: articlePublicH.CreateComment},
		{Method: http.MethodPost, Path: "/api/v1/articles/:id/like", Handler: articlePublicH.Like},
		{Method: http.MethodDelete, Path: "/api/v1/articles/:id/like", Handler: articlePublicH.Unlike},
		{Method: http.MethodPost, Path: "/api/v1/articles/:id/favorite", Handler: articlePublicH.Favorite},
		{Method: http.MethodDelete, Path: "/api/v1/articles/:id/favorite", Handler: articlePublicH.Unfavorite},
		{Method: http.MethodGet, Path: "/api/v1/articles/:id/engagement", Handler: articlePublicH.Status},
		{Method: http.MethodGet, Path: "/api/v1/user/article-favorites", Handler: articlePublicH.ListMyFavorites},
		{Method: http.MethodGet, Path: "/api/v1/user/article-likes", Handler: articlePublicH.ListMyLikes},
		{Method: http.MethodGet, Path: "/api/v1/user/articles", Handler: articlePublicH.ListMine},
		{Method: http.MethodPost, Path: "/api/v1/user/articles", Handler: articlePublicH.CreateMine},
		{Method: http.MethodGet, Path: "/api/v1/user/articles/:id", Handler: articlePublicH.DetailMine},
		{Method: http.MethodPut, Path: "/api/v1/user/articles/:id", Handler: articlePublicH.UpdateMine},
		{Method: http.MethodDelete, Path: "/api/v1/user/articles/:id", Handler: articlePublicH.DeleteMine},
		{Method: http.MethodPost, Path: "/api/v1/user/article-uploads", Handler: articlePublicH.UploadMine},
		{Method: http.MethodPost, Path: "/api/v1/user/favorites", Handler: favoriteUser.Add},
		{Method: http.MethodDelete, Path: "/api/v1/user/favorites/:product_id", Handler: favoriteUser.Remove},
		{Method: http.MethodPost, Path: "/api/v1/user/favorites/batch-remove", Handler: favoriteUser.RemoveBatch},
		{Method: http.MethodGet, Path: "/api/v1/user/favorites", Handler: favoriteUser.List},
		{Method: http.MethodGet, Path: "/api/v1/products/:id/favorite", Handler: favoriteUser.Status},
	}))
	server.AddRoutes(mws.MerchantOwner([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/merchant/products", Handler: adminH.List},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products", Handler: adminH.Create},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/batch", Handler: adminH.Batch},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/jobs/:id", Handler: adminH.JobStatus},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/recycle/restore", Handler: adminH.RecycleRestore},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/products/recycle", Handler: adminH.RecycleDelete},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/export", Handler: adminH.Export},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/import", Handler: adminH.Import},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/op-logs", Handler: adminH.OpLogs},
		{Method: http.MethodGet, Path: "/api/v1/merchant/products/:id", Handler: adminH.Detail},
		{Method: http.MethodPut, Path: "/api/v1/merchant/products/:id", Handler: adminH.Update},
		{Method: http.MethodPut, Path: "/api/v1/merchant/products/:id/status", Handler: adminH.SetStatus},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/:id/copy", Handler: adminH.Copy},
		{Method: http.MethodPost, Path: "/api/v1/merchant/products/:id/schedules", Handler: adminH.Schedule},
		{Method: http.MethodPut, Path: "/api/v1/merchant/skus/:id/stock", Handler: adminH.AdjustStock},
		{Method: http.MethodPost, Path: "/api/v1/merchant/skus/batch-stock", Handler: adminH.BatchStock},
		{Method: http.MethodGet, Path: "/api/v1/merchant/stocks/warnings", Handler: adminH.StockWarnings},
		{Method: http.MethodPost, Path: "/api/v1/merchant/uploads/images", Handler: adminH.Upload},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/schedules/:id", Handler: adminH.CancelSchedule},
		{Method: http.MethodGet, Path: "/api/v1/merchant/tags", Handler: adminH.ListTags},
		{Method: http.MethodPost, Path: "/api/v1/merchant/tags", Handler: adminH.SaveTag},
		{Method: http.MethodPut, Path: "/api/v1/merchant/tags/:id", Handler: adminH.SaveTag},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/tags/:id", Handler: adminH.DeleteTag},
		{Method: http.MethodGet, Path: "/api/v1/merchant/attr-templates", Handler: adminH.ListAttrTemplates},
		{Method: http.MethodPost, Path: "/api/v1/merchant/attr-templates", Handler: adminH.SaveAttrTemplate},
		{Method: http.MethodPut, Path: "/api/v1/merchant/attr-templates/:id", Handler: adminH.SaveAttrTemplate},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/attr-templates/:id", Handler: adminH.DeleteAttrTemplate},
		{Method: http.MethodGet, Path: "/api/v1/merchant/auth/me", Handler: shopOpsH.AuthMe},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/roles", Handler: shopOpsH.ListRoles},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/roles/:id/menus", Handler: shopOpsH.RoleMenus},
		{Method: http.MethodPost, Path: "/api/v1/merchant/shop/roles", Handler: shopOpsH.SaveRole},
		{Method: http.MethodPut, Path: "/api/v1/merchant/shop/roles/:id", Handler: shopOpsH.SaveRole},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/menus", Handler: shopOpsH.ListMenus},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shop/staff", Handler: shopOpsH.ListStaff},
		{Method: http.MethodPost, Path: "/api/v1/merchant/shop/staff", Handler: shopOpsH.BindStaff},
		{Method: http.MethodGet, Path: "/api/v1/merchant/articles", Handler: articleMerchantH.List},
		{Method: http.MethodPost, Path: "/api/v1/merchant/articles", Handler: articleMerchantH.Create},
		{Method: http.MethodGet, Path: "/api/v1/merchant/articles/:id", Handler: articleMerchantH.Detail},
		{Method: http.MethodPut, Path: "/api/v1/merchant/articles/:id", Handler: articleMerchantH.Update},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/articles/:id", Handler: articleMerchantH.Delete},
		{Method: http.MethodGet, Path: "/api/v1/merchant/article-categories", Handler: articleMerchantH.CategoryList},
		{Method: http.MethodGet, Path: "/api/v1/merchant/article-comments", Handler: articleMerchantH.CommentList},
		{Method: http.MethodPatch, Path: "/api/v1/merchant/article-comments/:id", Handler: articleMerchantH.CommentPatch},
		{Method: http.MethodDelete, Path: "/api/v1/merchant/article-comments/:id", Handler: articleMerchantH.CommentDelete},
		{Method: http.MethodPost, Path: "/api/v1/merchant/article-uploads", Handler: articleMerchantH.Upload},
		{Method: http.MethodGet, Path: "/api/v1/merchant/notifications/unread-count", Handler: notifH.UnreadCount},
		{Method: http.MethodPost, Path: "/api/v1/merchant/notifications/read-all", Handler: notifH.MarkAllRead},
		{Method: http.MethodGet, Path: "/api/v1/merchant/notifications", Handler: notifH.List},
		{Method: http.MethodPost, Path: "/api/v1/merchant/notifications/:id/read", Handler: notifH.MarkRead},
	}))
	server.AddRoutes(mws.PlatformAdmin([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/favorites", Handler: favoriteAdmin.AdminUserList},
		{Method: http.MethodGet, Path: "/api/v1/admin/products", Handler: platformProductH.List},
		{Method: http.MethodPut, Path: "/api/v1/admin/products/:id/off_sale", Handler: platformProductH.OffSale},
		{Method: http.MethodDelete, Path: "/api/v1/admin/products/:id", Handler: platformProductH.Delete},
		{Method: http.MethodGet, Path: "/api/v1/admin/categories", Handler: catalogAdmin.AdminListCategories},
		{Method: http.MethodPost, Path: "/api/v1/admin/categories", Handler: catalogAdmin.AdminCreateCategory},
		{Method: http.MethodPut, Path: "/api/v1/admin/categories/:id", Handler: catalogAdmin.AdminUpdateCategory},
		{Method: http.MethodDelete, Path: "/api/v1/admin/categories/:id", Handler: catalogAdmin.AdminDeleteCategory},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles/stats", Handler: articleAdminH.Stats},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles/recycle", Handler: articleAdminH.List},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/recycle/restore", Handler: articleAdminH.RecycleRestore},
		{Method: http.MethodDelete, Path: "/api/v1/admin/articles/recycle", Handler: articleAdminH.RecycleDelete},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/batch-audit", Handler: articleAdminH.BatchAudit},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles", Handler: articleAdminH.List},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles", Handler: articleAdminH.Create},
		{Method: http.MethodGet, Path: "/api/v1/admin/articles/:id", Handler: articleAdminH.Detail},
		{Method: http.MethodPut, Path: "/api/v1/admin/articles/:id", Handler: articleAdminH.Update},
		{Method: http.MethodDelete, Path: "/api/v1/admin/articles/:id", Handler: articleAdminH.SoftDelete},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/:id/audit", Handler: articleAdminH.Audit},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/:id/top", Handler: articleAdminH.Top},
		{Method: http.MethodPost, Path: "/api/v1/admin/articles/:id/offline", Handler: articleAdminH.Offline},
		{Method: http.MethodGet, Path: "/api/v1/admin/article-categories", Handler: articleAdminH.CategoryList},
		{Method: http.MethodPost, Path: "/api/v1/admin/article-categories", Handler: articleAdminH.CategoryCreate},
		{Method: http.MethodPut, Path: "/api/v1/admin/article-categories/:id", Handler: articleAdminH.CategoryUpdate},
		{Method: http.MethodDelete, Path: "/api/v1/admin/article-categories/:id", Handler: articleAdminH.CategoryDelete},
		{Method: http.MethodGet, Path: "/api/v1/admin/article-comments", Handler: articleAdminH.CommentList},
		{Method: http.MethodPatch, Path: "/api/v1/admin/article-comments/:id", Handler: articleAdminH.CommentPatch},
		{Method: http.MethodDelete, Path: "/api/v1/admin/article-comments/:id", Handler: articleAdminH.CommentDelete},
		{Method: http.MethodGet, Path: "/api/v1/admin/comment-emojis", Handler: articleAdminH.EmojiList},
		{Method: http.MethodPost, Path: "/api/v1/admin/comment-emojis", Handler: articleAdminH.EmojiCreate},
		{Method: http.MethodPut, Path: "/api/v1/admin/comment-emojis/:id", Handler: articleAdminH.EmojiUpdate},
		{Method: http.MethodDelete, Path: "/api/v1/admin/comment-emojis/:id", Handler: articleAdminH.EmojiDelete},
		{Method: http.MethodPost, Path: "/api/v1/admin/article-uploads", Handler: articleAdminH.Upload},
		{Method: http.MethodGet, Path: "/api/v1/admin/banners", Handler: articleAdminH.ListBanners},
		{Method: http.MethodPost, Path: "/api/v1/admin/banners", Handler: articleAdminH.CreateBanner},
		{Method: http.MethodPost, Path: "/api/v1/admin/banners/upload", Handler: articleAdminH.UploadBanner},
		{Method: http.MethodGet, Path: "/api/v1/admin/banners/:id", Handler: articleAdminH.GetBanner},
		{Method: http.MethodPut, Path: "/api/v1/admin/banners/:id", Handler: articleAdminH.UpdateBanner},
		{Method: http.MethodDelete, Path: "/api/v1/admin/banners/:id", Handler: articleAdminH.DeleteBanner},
		{Method: http.MethodPost, Path: "/api/v1/admin/shop-uploads", Handler: shopUploadH.Upload},
	}))
}

func serveProductsUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("products", httpserver.PathParam(r, "shop"), httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
func serveExportsUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("exports", httpserver.PathParam(r, "shop"), httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
func serveArticlesUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("articles", httpserver.PathParam(r, "shop"), httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
func serveShopsUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("shops", httpserver.PathParam(r, "owner"), httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
func serveReviewsUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("reviews", httpserver.PathParam(r, "user"), httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
func servePointsMallUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("points-mall", httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
func serveBannersUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("banners", httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
