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

type GetAdminRolesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetAdminRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminRolesLogic {
	return &GetAdminRolesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetAdminRolesLogic) GetAdminRoles(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	ids, err := biz.NewRBACLogic(l.svcCtx).AdminRoleIDs(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: ids}, nil
}
