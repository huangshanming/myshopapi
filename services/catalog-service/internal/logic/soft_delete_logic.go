package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type SoftDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSoftDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SoftDeleteLogic {
	return &SoftDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SoftDeleteLogic) SoftDelete(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).SoftDelete(w, r)
}
