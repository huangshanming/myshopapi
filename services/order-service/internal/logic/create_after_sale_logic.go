package logic

import (
	"net/http"

	"context"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAfterSaleLogic {
	return &CreateAfterSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAfterSaleLogic) CreateAfterSale(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).CreateAfterSale(w, r)
}
