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

type DeleteRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(ctx context.Context, req *types.IdPathReq) error {
	_, err := hadmin.NewAdminHandler(l.svcCtx).DeleteRole(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}})
	if err != nil {
		return err
	}
	return nil
}
