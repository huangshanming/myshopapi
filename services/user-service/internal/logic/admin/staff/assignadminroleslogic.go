// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package staff

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignAdminRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignAdminRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignAdminRolesLogic {
	return &AssignAdminRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignAdminRolesLogic) AssignAdminRoles(req *types.AdminRolesReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
