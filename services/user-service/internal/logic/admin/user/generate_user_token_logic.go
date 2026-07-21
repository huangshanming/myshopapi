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

type GenerateUserTokenLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGenerateUserTokenLogic(svcCtx *svc.ServiceContext) *GenerateUserTokenLogic {
	return &GenerateUserTokenLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *GenerateUserTokenLogic) GenerateUserToken(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/users/{Id}/token", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewAdminHandler(l.svcCtx).GenerateUserToken)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
