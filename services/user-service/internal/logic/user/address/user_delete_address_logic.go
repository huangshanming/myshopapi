package address

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserDeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteAddressLogic {
	return &UserDeleteAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteAddressLogic) UserDeleteAddress(w http.ResponseWriter, r *http.Request) {
	huser.NewAddressHandler(l.svcCtx).Delete(w, r)
}
