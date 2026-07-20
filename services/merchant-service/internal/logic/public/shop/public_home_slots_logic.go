package shop

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"
)

type PublicHomeSlotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicHomeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicHomeSlotsLogic {
	return &PublicHomeSlotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicHomeSlotsLogic) PublicHomeSlots(w http.ResponseWriter, r *http.Request) {
	hpublic.NewHomepageSlotHandler(l.svcCtx).PublicHomeSlots(w, r)
}
