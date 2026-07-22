package review

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

type MerchantDeleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantDeleteLogic {
	return &MerchantDeleteLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantDeleteLogic) MerchantDelete(ctx context.Context, req *types.IdPathReq) (*types.EmptyResp, error) {
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewReviewLogic(l.svcCtx).SoftDelete(ctx, req.Id, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
