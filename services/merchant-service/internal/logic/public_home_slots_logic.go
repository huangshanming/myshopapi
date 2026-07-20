package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
