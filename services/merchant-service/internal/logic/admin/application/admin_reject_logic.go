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

type AdminRejectLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRejectLogic {
	return &AdminRejectLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminRejectLogic) AdminReject(ctx context.Context, req *types.RejectBodyReq) (resp *types.EmptyResp, err error) {
	adminID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	appID := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).Reject(ctx, appID, adminID, req.Reason); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
