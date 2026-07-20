package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantRemarkLogic {
	return &MerchantRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantRemarkLogic) MerchantRemark(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewOrderHandler(l.svcCtx).MerchantRemark(w, r)
}
