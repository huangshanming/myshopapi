package role

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleMenusLogic {
	return &AssignRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleMenusLogic) AssignRoleMenus(ctx context.Context, req *types.RoleMenusReq) error {
	_, err := hadmin.NewAdminHandler(l.svcCtx).AssignRoleMenus(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Body: req})
	if err != nil {
		return err
	}
	return nil
}
