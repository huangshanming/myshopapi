package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type OpLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpLogsLogic {
	return &OpLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpLogsLogic) OpLogs(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).OpLogs(w, r)
}
