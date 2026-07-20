package internalapi

import (
	"net/http"

	"mymall/services/user-service/internal/logic/internalapi"
	"mymall/services/user-service/internal/svc"
)

func InternalCreateNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalCreateNotificationLogic(r.Context(), svcCtx)
		l.InternalCreateNotification(w, r)
	}
}

func InternalDeductPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalDeductPointsLogic(r.Context(), svcCtx)
		l.InternalDeductPoints(w, r)
	}
}

func InternalEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalEventLogic(r.Context(), svcCtx)
		l.InternalEvent(w, r)
	}
}

func InternalFreezeWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalFreezeWalletLogic(r.Context(), svcCtx)
		l.InternalFreezeWallet(w, r)
	}
}

func InternalGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalGetLogic(r.Context(), svcCtx)
		l.InternalGet(w, r)
	}
}

func InternalRefundPointsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalRefundPointsLogic(r.Context(), svcCtx)
		l.InternalRefundPoints(w, r)
	}
}

func InternalSettleWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalSettleWalletLogic(r.Context(), svcCtx)
		l.InternalSettleWallet(w, r)
	}
}

func InternalUnfreezeWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := internalapi.NewInternalUnfreezeWalletLogic(r.Context(), svcCtx)
		l.InternalUnfreezeWallet(w, r)
	}
}
