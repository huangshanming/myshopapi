package order

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAfterSaleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAfterSaleLogic {
	return &CreateAfterSaleLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *CreateAfterSaleLogic) CreateAfterSale(ctx context.Context, req *types.CreateAfterSaleBodyReq) (*types.AnyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	as, err := biz.NewOrderLogic(l.svcCtx).CreateAfterSale(ctx, userID, req.Id, types.CreateAfterSaleReq{
		Type: req.Type, Reason: req.Reason, Amount: req.Amount,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: as}, nil
}
