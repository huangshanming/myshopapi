package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type ListMyLikesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLikesLogic {
	return &ListMyLikesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyLikesLogic) ListMyLikes(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).ListMyLikes(w, r)
}
