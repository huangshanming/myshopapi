package shop

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MyShopsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMyShopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyShopsLogic {
	return &MyShopsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyShopsLogic) MyShops(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewShopHandler(l.svcCtx).MyShops(w, r)
}
