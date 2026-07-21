package role

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(ctx context.Context, req *types.RoleReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/roles", nil, nil, req, hadmin.NewAdminHandler(l.svcCtx).CreateRole)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
