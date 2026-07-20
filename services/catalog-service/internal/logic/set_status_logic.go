package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type SetStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetStatusLogic {
	return &SetStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetStatusLogic) SetStatus(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).SetStatus(w, r)
}
