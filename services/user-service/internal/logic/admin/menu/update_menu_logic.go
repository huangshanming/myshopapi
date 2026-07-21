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

type UpdateMenuLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(ctx context.Context, req *types.MenuUpdateReq) error {
	if err := biz.NewRBACLogic(l.svcCtx).UpdateMenu(ctx, req.Id, types.MenuReq{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Type:      req.Type,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Perms:     req.Perms,
		Sort:      req.Sort,
		Visible:   req.Visible,
		Status:    req.Status,
	}); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
