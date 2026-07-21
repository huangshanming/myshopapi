package user

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(ctx context.Context, req *types.UserUpdateReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/users/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).UpdateUser)
	if err != nil {
		return err
	}
	return nil
}
