package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type ListBannersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBannersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBannersLogic {
	return &ListBannersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBannersLogic) ListBanners(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).ListBanners(w, r)
}
