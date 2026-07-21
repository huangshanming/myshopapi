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

type AssignRoleMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignRoleMenusLogic(svcCtx *svc.ServiceContext) *AssignRoleMenusLogic {
	return &AssignRoleMenusLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleMenusLogic) AssignRoleMenus(ctx context.Context, req *types.RoleMenusReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/roles/{Id}/menus", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).AssignRoleMenus)
	if err != nil {
		return err
	}
	return nil
}
