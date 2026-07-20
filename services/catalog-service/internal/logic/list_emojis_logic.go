package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type ListEmojisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmojisLogic {
	return &ListEmojisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEmojisLogic) ListEmojis(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).ListEmojis(w, r)
}
