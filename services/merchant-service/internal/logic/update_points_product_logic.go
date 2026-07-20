package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePointsProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePointsProductLogic {
	return &UpdatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePointsProductLogic) UpdatePointsProduct(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsProductHandler(l.svcCtx).Update(w, r)
}
