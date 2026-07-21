package shopops

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hhandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRoleMenusLogic(svcCtx *svc.ServiceContext) *RoleMenusLogic {
	return &RoleMenusLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *RoleMenusLogic) RoleMenus(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/shop/roles/:id/menus", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hhandler.NewShopOpsHandler(l.svcCtx).RoleMenus)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
