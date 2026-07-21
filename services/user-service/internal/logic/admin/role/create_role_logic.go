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

type CreateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(ctx context.Context, req *types.RoleReq) (resp *types.AnyResp, err error) {
	role, err := biz.NewRBACLogic(l.svcCtx).CreateRole(ctx, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: role}, nil
}
