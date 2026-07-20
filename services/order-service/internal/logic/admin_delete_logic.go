package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteLogic {
	return &AdminDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteLogic) AdminDelete(w http.ResponseWriter, r *http.Request) {
	hadmin.NewReviewHandler(l.svcCtx).AdminDelete(w, r)
}
