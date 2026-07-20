package handler

import (
	"net/http"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

func AssignRoleMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAssignRoleMenusLogic(r.Context(), svcCtx)
		l.AssignRoleMenus(w, r)
	}
}
