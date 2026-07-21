package role

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetRoleMenusLogic(svcCtx *svc.ServiceContext) *GetRoleMenusLogic {
	return &GetRoleMenusLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *GetRoleMenusLogic) GetRoleMenus(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/roles/{Id}/menus", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewAdminHandler(l.svcCtx).GetRoleMenus)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
