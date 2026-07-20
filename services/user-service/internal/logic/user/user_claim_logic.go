package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserClaimLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClaimLogic {
	return &UserClaimLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClaimLogic) UserClaim(w http.ResponseWriter, r *http.Request) {
	huser.NewTaskHandler(l.svcCtx).UserClaim(w, r)
}
