package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
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
