package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type BatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchLogic {
	return &BatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchLogic) Batch(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).Batch(w, r)
}
