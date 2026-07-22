package shop

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateShopLogic {
	return &AdminUpdateShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateShopLogic) AdminUpdateShop(ctx context.Context, req *types.AdminUpdateShopBodyReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).AdminUpdateShop(ctx, id, req.ToAdminUpdateShopReq()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
