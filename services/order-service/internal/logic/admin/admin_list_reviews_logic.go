package admin

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"
)

type AdminListReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListReviewsLogic {
	return &AdminListReviewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListReviewsLogic) AdminListReviews(w http.ResponseWriter, r *http.Request) {
	hadmin.NewReviewHandler(l.svcCtx).AdminList(w, r)
}
