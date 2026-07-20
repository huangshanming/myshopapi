package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserPointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPointsLogic {
	return &UserPointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPointsLogic) UserPoints(w http.ResponseWriter, r *http.Request) {
	huser.NewTaskHandler(l.svcCtx).UserPoints(w, r)
}
