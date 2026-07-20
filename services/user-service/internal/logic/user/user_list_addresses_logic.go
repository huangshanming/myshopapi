package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserListAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListAddressesLogic {
	return &UserListAddressesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListAddressesLogic) UserListAddresses(w http.ResponseWriter, r *http.Request) {
	huser.NewAddressHandler(l.svcCtx).Create(w, r)
}
