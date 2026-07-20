package logic

import (
	"net/http"

	"context"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EligibleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEligibleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EligibleLogic {
	return &EligibleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EligibleLogic) Eligible(w http.ResponseWriter, r *http.Request) {
	huser.NewReviewHandler(l.svcCtx).Eligible(w, r)
}
