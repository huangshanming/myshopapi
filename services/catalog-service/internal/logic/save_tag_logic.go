package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type SaveTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveTagLogic {
	return &SaveTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveTagLogic) SaveTag(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).SaveTag(w, r)
}
