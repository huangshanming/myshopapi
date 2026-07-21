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

type DeleteRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(ctx context.Context, req *types.IdPathReq) error {
	_, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/admin/roles/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewAdminHandler(l.svcCtx).DeleteRole)
	if err != nil {
		return err
	}
	return nil
}
