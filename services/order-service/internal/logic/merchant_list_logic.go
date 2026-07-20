package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListLogic {
	return &MerchantListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListLogic) MerchantList(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewReviewHandler(l.svcCtx).MerchantList(w, r)
}
