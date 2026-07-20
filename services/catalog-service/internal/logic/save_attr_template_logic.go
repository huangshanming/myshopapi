package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type SaveAttrTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAttrTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAttrTemplateLogic {
	return &SaveAttrTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveAttrTemplateLogic) SaveAttrTemplate(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).SaveAttrTemplate(w, r)
}
