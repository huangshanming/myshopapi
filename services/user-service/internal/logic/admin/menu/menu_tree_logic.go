package menu

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuTreeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMenuTreeLogic(svcCtx *svc.ServiceContext) *MenuTreeLogic {
	return &MenuTreeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MenuTreeLogic) MenuTree(ctx context.Context) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/menus", nil, nil, nil, hadmin.NewAdminHandler(l.svcCtx).MenuTree)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
