package notification

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkNotificationReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMarkNotificationReadLogic(svcCtx *svc.ServiceContext) *MarkNotificationReadLogic {
	return &MarkNotificationReadLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MarkNotificationReadLogic) MarkNotificationRead(ctx context.Context, req *types.IdPathReq) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/notifications/{Id}/read", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, huser.NewUserHandler(l.svcCtx).MarkNotificationRead)
	if err != nil {
		return err
	}
	return nil
}
