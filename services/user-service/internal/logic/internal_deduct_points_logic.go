package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalDeductPointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalDeductPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalDeductPointsLogic {
	return &InternalDeductPointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalDeductPointsLogic) InternalDeductPoints(w http.ResponseWriter, r *http.Request) {
	hinternal.NewTaskHandler(l.svcCtx).InternalDeductPoints(w, r)
}
