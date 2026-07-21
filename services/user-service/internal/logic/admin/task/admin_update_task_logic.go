package task

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateTaskLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateTaskLogic(svcCtx *svc.ServiceContext) *AdminUpdateTaskLogic {
	return &AdminUpdateTaskLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateTaskLogic) AdminUpdateTask(ctx context.Context, req *types.UpdateTaskReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/tasks/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewTaskHandler(l.svcCtx).AdminUpdate)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
