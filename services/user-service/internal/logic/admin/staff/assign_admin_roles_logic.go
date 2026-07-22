package staff

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AssignAdminRolesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignAdminRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignAdminRolesLogic {
	return &AssignAdminRolesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AssignAdminRolesLogic) AssignAdminRoles(ctx context.Context, req *types.AdminRolesReq) (*types.EmptyResp, error) {
	if err := biz.NewRBACLogic(l.svcCtx).AssignAdminRoles(ctx, req.Id, req.RoleIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
