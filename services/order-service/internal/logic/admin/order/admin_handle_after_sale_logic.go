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

type AdminHandleAfterSaleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminHandleAfterSaleLogic {
	return &AdminHandleAfterSaleLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminHandleAfterSaleLogic) AdminHandleAfterSale(ctx context.Context, req *types.HandleAfterSaleBodyReq) (*types.EmptyResp, error) {
	uid, _ := middleware.GetUserID(ctx)
	if err := biz.NewOrderLogic(l.svcCtx).HandleAfterSale(ctx, req.Id, 0, uid, req.Action, req.AdminRemark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
