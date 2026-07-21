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

type AdminDisableShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDisableShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDisableShopLogic {
	return &AdminDisableShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminDisableShopLogic) AdminDisableShop(ctx context.Context, req *types.RejectBodyReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).DisableShop(ctx, id, req.Reason); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
