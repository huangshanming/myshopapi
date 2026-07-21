package menu

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(ctx context.Context, req *types.MenuUpdateReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/menus/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).UpdateMenu)
	if err != nil {
		return err
	}
	return nil
}
