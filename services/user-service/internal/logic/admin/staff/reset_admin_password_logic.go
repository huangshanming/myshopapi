package staff

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAdminPasswordLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewResetAdminPasswordLogic(svcCtx *svc.ServiceContext) *ResetAdminPasswordLogic {
	return &ResetAdminPasswordLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ResetAdminPasswordLogic) ResetAdminPassword(ctx context.Context, req *types.AdminResetPwdReq) error {
	_, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/admins/{Id}/password", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewAdminHandler(l.svcCtx).ResetAdminPassword)
	if err != nil {
		return err
	}
	return nil
}
