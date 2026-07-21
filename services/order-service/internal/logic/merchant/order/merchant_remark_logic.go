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

type MerchantRemarkLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantRemarkLogic {
	return &MerchantRemarkLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantRemarkLogic) MerchantRemark(ctx context.Context, req *types.RemarkBodyReq) (*types.EmptyResp, error) {
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewOrderLogic(l.svcCtx).UpdateRemark(ctx, req.Id, shopID, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
