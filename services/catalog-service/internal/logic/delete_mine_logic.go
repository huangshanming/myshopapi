package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type DeleteMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMineLogic {
	return &DeleteMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMineLogic) DeleteMine(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).DeleteMine(w, r)
}
