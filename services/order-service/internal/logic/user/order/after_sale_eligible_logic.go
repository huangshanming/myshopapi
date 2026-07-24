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

type AfterSaleEligibleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAfterSaleEligibleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterSaleEligibleLogic {
	return &AfterSaleEligibleLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AfterSaleEligibleLogic) AfterSaleEligible(ctx context.Context, req *types.IdPathReq) (*types.AfterSaleEligibleResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	resp, err := biz.NewOrderLogic(l.svcCtx).AfterSaleEligible(ctx, userID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return resp, nil
}
