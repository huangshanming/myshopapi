package homepage

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

type AdminGrantSlotLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGrantSlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantSlotLogic {
	return &AdminGrantSlotLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantSlotLogic) AdminGrantSlot(ctx context.Context, req *types.GrantSlotReq) (resp *types.HomepageOrderResp, err error) {
	adminID, _ := middleware.GetUserID(ctx)
	order, err := biz.NewMerchantLogic(l.svcCtx).GrantSlot(adminID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.HomepageOrderResp{Data: order}, nil
}
