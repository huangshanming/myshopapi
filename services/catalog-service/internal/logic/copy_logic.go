package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type CopyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCopyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CopyLogic {
	return &CopyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CopyLogic) Copy(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).Copy(w, r)
}
