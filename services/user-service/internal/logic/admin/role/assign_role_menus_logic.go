package role

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
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
	if err := biz.NewRBACLogic(l.svcCtx).AssignRoleMenus(ctx, req.Id, req.MenuIDs); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
