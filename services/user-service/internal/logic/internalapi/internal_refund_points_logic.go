package internalapi

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"
)

type InternalRefundPointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalRefundPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalRefundPointsLogic {
	return &InternalRefundPointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalRefundPointsLogic) InternalRefundPoints(w http.ResponseWriter, r *http.Request) {
	hinternal.NewTaskHandler(l.svcCtx).InternalRefundPoints(w, r)
}
