package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type UpdateMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMineLogic {
	return &UpdateMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMineLogic) UpdateMine(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).UpdateMine(w, r)
}
