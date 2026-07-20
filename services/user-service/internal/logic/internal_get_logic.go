package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalGetLogic {
	return &InternalGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalGetLogic) InternalGet(w http.ResponseWriter, r *http.Request) {
	hinternal.NewAddressHandler(l.svcCtx).InternalGet(w, r)
}
