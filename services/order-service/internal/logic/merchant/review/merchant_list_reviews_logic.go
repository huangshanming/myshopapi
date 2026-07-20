package review

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"
)

type MerchantListReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListReviewsLogic {
	return &MerchantListReviewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListReviewsLogic) MerchantListReviews(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewReviewHandler(l.svcCtx).MerchantList(w, r)
}
