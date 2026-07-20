package address

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/address"
	"mymall/services/user-service/internal/svc"
)

func SetDefaultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewSetDefaultLogic(r.Context(), svcCtx)
		l.SetDefault(w, r)
	}
}

func UserCreateAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewUserCreateAddressLogic(r.Context(), svcCtx)
		l.UserCreateAddress(w, r)
	}
}

func UserDeleteAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewUserDeleteAddressLogic(r.Context(), svcCtx)
		l.UserDeleteAddress(w, r)
	}
}

func UserListAddressesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewUserListAddressesLogic(r.Context(), svcCtx)
		l.UserListAddresses(w, r)
	}
}

func UserUpdateAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewUserUpdateAddressLogic(r.Context(), svcCtx)
		l.UserUpdateAddress(w, r)
	}
}
