package logic

import (
	"net/http"

	"context"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultLogic {
	return &SetDefaultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultLogic) SetDefault(w http.ResponseWriter, r *http.Request) {
	huser.NewAddressHandler(l.svcCtx).SetDefault(w, r)
}
