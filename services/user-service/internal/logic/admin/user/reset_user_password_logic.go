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

type ResetUserPasswordLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewResetUserPasswordLogic(svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(ctx context.Context, req *types.UserResetPwdReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/users/{Id}/password", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).ResetUserPassword)
	if err != nil {
		return err
	}
	return nil
}
