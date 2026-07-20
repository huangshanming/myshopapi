package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type RecycleRestoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecycleRestoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleRestoreLogic {
	return &RecycleRestoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleRestoreLogic) RecycleRestore(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).RecycleRestore(w, r)
}
