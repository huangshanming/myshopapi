package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type BatchAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchAuditLogic {
	return &BatchAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchAuditLogic) BatchAudit(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).BatchAudit(w, r)
}
