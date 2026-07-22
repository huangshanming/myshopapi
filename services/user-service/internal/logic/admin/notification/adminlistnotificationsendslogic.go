// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListNotificationSendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListNotificationSendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListNotificationSendsLogic {
	return &AdminListNotificationSendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListNotificationSendsLogic) AdminListNotificationSends(req *types.PageReq) (resp *types.PageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
