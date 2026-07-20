package review

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type GetByOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetByOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByOrderLogic {
	return &GetByOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetByOrderLogic) GetByOrder(w http.ResponseWriter, r *http.Request) {
	huser.NewReviewHandler(l.svcCtx).GetByOrder(w, r)
}
