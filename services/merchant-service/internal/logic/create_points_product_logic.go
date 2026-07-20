package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePointsProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePointsProductLogic {
	return &CreatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePointsProductLogic) CreatePointsProduct(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsProductHandler(l.svcCtx).Create(w, r)
}
