package merchant

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantListSeckillEntriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListSeckillEntriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListSeckillEntriesLogic {
	return &MerchantListSeckillEntriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListSeckillEntriesLogic) MerchantListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewSeckillHandler(l.svcCtx).MerchantListSeckillEntries(w, r)
}
