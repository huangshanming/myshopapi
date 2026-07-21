package menu

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type MenuTreeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuTreeLogic {
	return &MenuTreeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MenuTreeLogic) MenuTree(ctx context.Context) (resp *types.AnyResp, err error) {
	tree, err := biz.NewRBACLogic(l.svcCtx).MenuTreeAll(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: tree}, nil
}
