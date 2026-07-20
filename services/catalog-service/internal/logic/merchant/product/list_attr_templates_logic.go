package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type ListAttrTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAttrTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAttrTemplatesLogic {
	return &ListAttrTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAttrTemplatesLogic) ListAttrTemplates(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).ListAttrTemplates(w, r)
}
