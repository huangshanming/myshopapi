package address

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserUpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateAddressLogic {
	return &UserUpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateAddressLogic) UserUpdateAddress(w http.ResponseWriter, r *http.Request) {
	huser.NewAddressHandler(l.svcCtx).Delete(w, r)
}
