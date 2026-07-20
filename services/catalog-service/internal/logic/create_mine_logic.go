package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type CreateMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMineLogic {
	return &CreateMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMineLogic) CreateMine(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).CreateMine(w, r)
}
