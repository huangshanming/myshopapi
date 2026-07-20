package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserCheckinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCheckinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCheckinLogic {
	return &UserCheckinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCheckinLogic) UserCheckin(w http.ResponseWriter, r *http.Request) {
	huser.NewTaskHandler(l.svcCtx).UserCheckin(w, r)
}
