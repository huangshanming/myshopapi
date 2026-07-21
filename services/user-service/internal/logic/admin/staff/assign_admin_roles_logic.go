package staff

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignAdminRolesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignAdminRolesLogic(svcCtx *svc.ServiceContext) *AssignAdminRolesLogic {
	return &AssignAdminRolesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AssignAdminRolesLogic) AssignAdminRoles(ctx context.Context, req *types.AdminRolesReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/admins/{Id}/roles", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).AssignAdminRoles)
	if err != nil {
		return err
	}
	return nil
}
