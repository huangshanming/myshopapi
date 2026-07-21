package application

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/merchant-service/internal/app/admin"
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
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewShopHandler(l.svcCtx).AdminReject(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
