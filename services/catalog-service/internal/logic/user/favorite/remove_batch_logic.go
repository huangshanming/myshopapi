package favorite

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type RemoveBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBatchLogic {
	return &RemoveBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveBatchLogic) RemoveBatch(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).RemoveBatch(w, r)
}
