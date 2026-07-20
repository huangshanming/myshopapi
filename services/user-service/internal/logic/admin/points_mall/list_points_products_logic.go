package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type ListPointsProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPointsProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointsProductsLogic {
	return &ListPointsProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPointsProductsLogic) ListPointsProducts(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsProductHandler(l.svcCtx).List(w, r)
}
