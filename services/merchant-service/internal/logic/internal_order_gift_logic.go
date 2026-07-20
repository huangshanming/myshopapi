package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalOrderGiftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalOrderGiftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalOrderGiftLogic {
	return &InternalOrderGiftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalOrderGiftLogic) InternalOrderGift(w http.ResponseWriter, r *http.Request) {
	hinternal.NewCouponHandler(l.svcCtx).InternalOrderGift(w, r)
}
