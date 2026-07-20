package article

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/user/article"
	"mymall/services/catalog-service/internal/svc"
)

func CreateCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewCreateCommentLogic(r.Context(), svcCtx)
		l.CreateComment(w, r)
	}
}

func CreateMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewCreateMineLogic(r.Context(), svcCtx)
		l.CreateMine(w, r)
	}
}

func DeleteMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewDeleteMineLogic(r.Context(), svcCtx)
		l.DeleteMine(w, r)
	}
}

func DetailMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewDetailMineLogic(r.Context(), svcCtx)
		l.DetailMine(w, r)
	}
}

func FavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewFavoriteLogic(r.Context(), svcCtx)
		l.Favorite(w, r)
	}
}

func LikeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewLikeLogic(r.Context(), svcCtx)
		l.Like(w, r)
	}
}

func ListMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewListMineLogic(r.Context(), svcCtx)
		l.ListMine(w, r)
	}
}

func ListMyFavoritesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewListMyFavoritesLogic(r.Context(), svcCtx)
		l.ListMyFavorites(w, r)
	}
}

func ListMyLikesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewListMyLikesLogic(r.Context(), svcCtx)
		l.ListMyLikes(w, r)
	}
}

func UnfavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewUnfavoriteLogic(r.Context(), svcCtx)
		l.Unfavorite(w, r)
	}
}

func UnlikeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewUnlikeLogic(r.Context(), svcCtx)
		l.Unlike(w, r)
	}
}

func UpdateMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewUpdateMineLogic(r.Context(), svcCtx)
		l.UpdateMine(w, r)
	}
}

func UploadMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewUploadMineLogic(r.Context(), svcCtx)
		l.UploadMine(w, r)
	}
}

func UserArticleEngagementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewUserArticleEngagementLogic(r.Context(), svcCtx)
		l.UserArticleEngagement(w, r)
	}
}
