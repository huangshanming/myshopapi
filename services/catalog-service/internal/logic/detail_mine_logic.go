package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type DetailMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailMineLogic {
	return &DetailMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailMineLogic) DetailMine(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).DetailMine(w, r)
}
