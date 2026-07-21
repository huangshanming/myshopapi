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

type DeleteMenuLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(ctx context.Context, req *types.IdPathReq) error {
	_, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/admin/menus/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewAdminHandler(l.svcCtx).DeleteMenu)
	if err != nil {
		return err
	}
	return nil
}
