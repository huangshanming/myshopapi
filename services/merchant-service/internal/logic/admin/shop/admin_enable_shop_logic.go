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

type AdminEnableShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminEnableShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminEnableShopLogic {
	return &AdminEnableShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminEnableShopLogic) AdminEnableShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).EnableShop(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
