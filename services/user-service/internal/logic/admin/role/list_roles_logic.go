package role

import (
	"context"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListRolesLogic) ListRoles(ctx context.Context) (resp *types.PageListResp, err error) {
	data, err := hadmin.NewAdminHandler(l.svcCtx).ListRoles(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.PageListResp{List: data}, nil
}
