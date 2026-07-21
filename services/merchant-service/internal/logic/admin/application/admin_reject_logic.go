package application

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hadmin "mymall/services/merchant-service/internal/app/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRejectLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminRejectLogic(svcCtx *svc.ServiceContext) *AdminRejectLogic {
	return &AdminRejectLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminRejectLogic) AdminReject(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/applications/:id/reject", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewShopHandler(l.svcCtx).AdminReject)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
