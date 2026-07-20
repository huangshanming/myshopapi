package banner

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/banner"
	"mymall/services/catalog-service/internal/svc"
)

func AdminListBannersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner.NewAdminListBannersLogic(r.Context(), svcCtx)
		l.AdminListBanners(w, r)
	}
}

func CreateBannerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner.NewCreateBannerLogic(r.Context(), svcCtx)
		l.CreateBanner(w, r)
	}
}

func DeleteBannerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner.NewDeleteBannerLogic(r.Context(), svcCtx)
		l.DeleteBanner(w, r)
	}
}

func GetBannerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner.NewGetBannerLogic(r.Context(), svcCtx)
		l.GetBanner(w, r)
	}
}

func UpdateBannerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner.NewUpdateBannerLogic(r.Context(), svcCtx)
		l.UpdateBanner(w, r)
	}
}

func UploadBannerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner.NewUploadBannerLogic(r.Context(), svcCtx)
		l.UploadBanner(w, r)
	}
}
