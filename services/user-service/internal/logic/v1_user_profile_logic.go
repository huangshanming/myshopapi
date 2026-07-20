package logic

import (
	"net/http"

	"context"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type V1UserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewV1UserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *V1UserProfileLogic {
	return &V1UserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *V1UserProfileLogic) V1UserProfile(w http.ResponseWriter, r *http.Request) {
	huser.NewUserHandler(l.svcCtx).Profile(w, r)
}
