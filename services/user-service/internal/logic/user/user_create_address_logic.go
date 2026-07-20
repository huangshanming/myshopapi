package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserCreateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateAddressLogic {
	return &UserCreateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateAddressLogic) UserCreateAddress(w http.ResponseWriter, r *http.Request) {
	huser.NewAddressHandler(l.svcCtx).Create(w, r)
}
