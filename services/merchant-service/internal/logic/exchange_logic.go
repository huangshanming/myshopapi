package logic

import (
	"net/http"

	"context"

	huser "mymall/services/merchant-service/internal/httpapi/user"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeLogic) Exchange(w http.ResponseWriter, r *http.Request) {
	huser.NewPointsOrderHandler(l.svcCtx).Exchange(w, r)
}
