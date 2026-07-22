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

type GetRoleMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleMenusLogic {
	return &GetRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetRoleMenusLogic) GetRoleMenus(ctx context.Context, req *types.IdPathReq) (resp *types.IdListResp, err error) {
	ids, err := biz.NewRBACLogic(l.svcCtx).GetRoleMenus(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.IdListResp{Ids: ids}, nil
}
