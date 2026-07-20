package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type OfflineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OfflineLogic {
	return &OfflineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OfflineLogic) Offline(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).Offline(w, r)
}
