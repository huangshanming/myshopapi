package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantList2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantList2Logic {
	return &MerchantList2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantList2Logic) MerchantList2(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewOrderHandler(l.svcCtx).MerchantList(w, r)
}
