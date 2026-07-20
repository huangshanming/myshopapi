package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type UserArticleEngagementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserArticleEngagementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserArticleEngagementLogic {
	return &UserArticleEngagementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserArticleEngagementLogic) UserArticleEngagement(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).Status(w, r)
}
