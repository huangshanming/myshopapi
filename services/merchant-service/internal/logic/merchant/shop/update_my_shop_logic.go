package shop

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMyShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyShopLogic {
	return &UpdateMyShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyShopLogic) UpdateMyShop(ctx context.Context, req *types.UpdateShopBodyReq) (resp *types.EmptyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	shopID := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateMyShop(ctx, shopID, userID, req.ToUpdateShopReq()); err != nil {
		return nil, xerr.New(http.StatusForbidden, err.Error())
	}
	return &types.EmptyResp{}, nil
}
