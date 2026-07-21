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

type GetUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/users/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewAdminHandler(l.svcCtx).GetUser)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
