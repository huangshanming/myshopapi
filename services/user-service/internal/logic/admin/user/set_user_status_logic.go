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

type SetUserStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSetUserStatusLogic(svcCtx *svc.ServiceContext) *SetUserStatusLogic {
	return &SetUserStatusLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *SetUserStatusLogic) SetUserStatus(ctx context.Context, req *types.UserStatusReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/users/{Id}/status", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).SetUserStatus)
	if err != nil {
		return err
	}
	return nil
}
