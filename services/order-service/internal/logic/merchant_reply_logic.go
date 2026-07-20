package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantReplyLogic {
	return &MerchantReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantReplyLogic) MerchantReply(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewReviewHandler(l.svcCtx).MerchantReply(w, r)
}
