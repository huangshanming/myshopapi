package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type UploadBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadBannerLogic {
	return &UploadBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadBannerLogic) UploadBanner(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).UploadBanner(w, r)
}
