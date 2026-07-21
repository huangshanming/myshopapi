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

type GetAdminRolesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetAdminRolesLogic(svcCtx *svc.ServiceContext) *GetAdminRolesLogic {
	return &GetAdminRolesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *GetAdminRolesLogic) GetAdminRoles(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/admins/{Id}/roles", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewAdminHandler(l.svcCtx).GetAdminRoles)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
