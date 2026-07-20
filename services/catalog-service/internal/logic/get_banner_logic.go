package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type GetBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBannerLogic {
	return &GetBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBannerLogic) GetBanner(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).GetBanner(w, r)
}
