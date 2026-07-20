package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type DeleteAttrTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAttrTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAttrTemplateLogic {
	return &DeleteAttrTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAttrTemplateLogic) DeleteAttrTemplate(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).DeleteAttrTemplate(w, r)
}
