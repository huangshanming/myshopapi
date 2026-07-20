package menu

import (
	"net/http"

	"mymall/services/user-service/internal/logic/admin/menu"
	"mymall/services/user-service/internal/svc"
)

func CreateMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewCreateMenuLogic(r.Context(), svcCtx)
		l.CreateMenu(w, r)
	}
}

func DeleteMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewDeleteMenuLogic(r.Context(), svcCtx)
		l.DeleteMenu(w, r)
	}
}

func MenuTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewMenuTreeLogic(r.Context(), svcCtx)
		l.MenuTree(w, r)
	}
}

func UpdateMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewUpdateMenuLogic(r.Context(), svcCtx)
		l.UpdateMenu(w, r)
	}
}
