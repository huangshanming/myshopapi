package comment

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/comment"
	"mymall/services/catalog-service/internal/svc"
)

func AdminDeleteArticleCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewAdminDeleteArticleCommentLogic(r.Context(), svcCtx)
		l.AdminDeleteArticleComment(w, r)
	}
}

func AdminListArticleCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewAdminListArticleCommentsLogic(r.Context(), svcCtx)
		l.AdminListArticleComments(w, r)
	}
}

func AdminPatchArticleCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewAdminPatchArticleCommentLogic(r.Context(), svcCtx)
		l.AdminPatchArticleComment(w, r)
	}
}

func EmojiCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewEmojiCreateLogic(r.Context(), svcCtx)
		l.EmojiCreate(w, r)
	}
}

func EmojiDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewEmojiDeleteLogic(r.Context(), svcCtx)
		l.EmojiDelete(w, r)
	}
}

func EmojiListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewEmojiListLogic(r.Context(), svcCtx)
		l.EmojiList(w, r)
	}
}

func EmojiUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := comment.NewEmojiUpdateLogic(r.Context(), svcCtx)
		l.EmojiUpdate(w, r)
	}
}
