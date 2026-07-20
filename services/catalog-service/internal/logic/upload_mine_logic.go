package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type UploadMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMineLogic {
	return &UploadMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadMineLogic) UploadMine(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).UploadMine(w, r)
}
