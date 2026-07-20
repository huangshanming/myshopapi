package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicListShopsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicListShopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListShopsLogic {
	return &PublicListShopsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicListShopsLogic) PublicListShops(w http.ResponseWriter, r *http.Request) {
	hpublic.NewShopHandler(l.svcCtx).PublicListShops(w, r)
}
