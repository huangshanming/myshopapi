package task

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"
)

type InternalEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalEventLogic {
	return &InternalEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalEventLogic) InternalEvent(w http.ResponseWriter, r *http.Request) {
	hinternal.NewTaskHandler(l.svcCtx).InternalEvent(w, r)
}
