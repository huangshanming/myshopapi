package merchant

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type UpdateMyShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyShopLogic {
	return &UpdateMyShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyShopLogic) UpdateMyShop(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewShopHandler(l.svcCtx).UpdateMyShop(w, r)
}
