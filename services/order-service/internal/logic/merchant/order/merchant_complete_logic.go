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

type MerchantCompleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCompleteLogic {
	return &MerchantCompleteLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantCompleteLogic) MerchantComplete(ctx context.Context, req *types.IdPathReq) (*types.EmptyResp, error) {
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewOrderLogic(l.svcCtx).Complete(ctx, req.Id, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
