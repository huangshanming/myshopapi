package shop

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

type AdminEnableShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminEnableShopLogic(svcCtx *svc.ServiceContext) *AdminEnableShopLogic {
	return &AdminEnableShopLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminEnableShopLogic) AdminEnableShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/shops/:id/enable", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewShopHandler(l.svcCtx).AdminEnableShop)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
