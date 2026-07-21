package task

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateTaskLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateTaskLogic {
	return &AdminUpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateTaskLogic) AdminUpdateTask(ctx context.Context, req *types.UpdateTaskReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewTaskHandler(l.svcCtx).AdminUpdate(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
