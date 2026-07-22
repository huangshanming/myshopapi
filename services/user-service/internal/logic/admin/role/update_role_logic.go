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

type UpdateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(ctx context.Context, req *types.RoleUpdateReq) (*types.EmptyResp, error) {
	if err := biz.NewRBACLogic(l.svcCtx).UpdateRole(ctx, req.Id, types.RoleReq{
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
		Remark: req.Remark,
	}); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
