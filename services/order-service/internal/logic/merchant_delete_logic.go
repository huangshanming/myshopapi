package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantDeleteLogic {
	return &MerchantDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantDeleteLogic) MerchantDelete(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewReviewHandler(l.svcCtx).MerchantDelete(w, r)
}
