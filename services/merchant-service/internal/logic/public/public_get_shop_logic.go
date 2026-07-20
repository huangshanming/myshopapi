package public

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"
)

type PublicGetShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicGetShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicGetShopLogic {
	return &PublicGetShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicGetShopLogic) PublicGetShop(w http.ResponseWriter, r *http.Request) {
	hpublic.NewShopHandler(l.svcCtx).PublicGetShop(w, r)
}
