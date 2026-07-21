package application

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

type AdminApproveLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminApproveLogic {
	return &AdminApproveLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminApproveLogic) AdminApprove(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	adminID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	appID := req.Id
	shop, err := biz.NewMerchantLogic(l.svcCtx).Approve(ctx, appID, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: shop}, nil
}
