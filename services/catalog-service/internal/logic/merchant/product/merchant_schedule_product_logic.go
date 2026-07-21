package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantScheduleProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantScheduleProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantScheduleProductLogic {
	return &MerchantScheduleProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantScheduleProductLogic) MerchantScheduleProduct(ctx context.Context, req *types.ScheduleBodyReq) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id := req.Id
	if err := plogic.NewProductAdminLogic(l.svcCtx).CreateSchedule(ctx, shopID, uid, id, req.ToProduct()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
