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

type UpdateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(ctx context.Context, req *types.RoleUpdateReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/roles/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).UpdateRole)
	if err != nil {
		return err
	}
	return nil
}
