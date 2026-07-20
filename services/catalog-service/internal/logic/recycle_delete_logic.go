package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type RecycleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecycleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleDeleteLogic {
	return &RecycleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleDeleteLogic) RecycleDelete(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).RecycleDelete(w, r)
}
