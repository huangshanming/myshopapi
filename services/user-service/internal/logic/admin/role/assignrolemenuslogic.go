// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleMenusLogic {
	return &AssignRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleMenusLogic) AssignRoleMenus(req *types.RoleMenusReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
