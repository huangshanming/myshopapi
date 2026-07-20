package handler

import (
	"net/http"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

func EligibleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewEligibleLogic(r.Context(), svcCtx)
		l.Eligible(w, r)
	}
}
