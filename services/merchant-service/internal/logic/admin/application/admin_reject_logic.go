package application

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

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

func (l *AdminRejectLogic) AdminReject(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	adminID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	appID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "申请ID无效")
	}
	var body types.RejectReq
	_ = appinput.BindBody(in, &body)
	if err := biz.NewMerchantLogic(l.svcCtx).Reject(ctx, appID, adminID, body.Reason); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
