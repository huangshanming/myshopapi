package staff

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(ctx context.Context, req *types.AdminCreateReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/admins", nil, nil, req, hadmin.NewAdminHandler(l.svcCtx).CreateAdmin)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
