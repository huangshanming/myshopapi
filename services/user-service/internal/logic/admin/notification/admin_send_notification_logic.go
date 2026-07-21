package notification

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AdminSendNotificationLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminSendNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminSendNotificationLogic {
	return &AdminSendNotificationLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminSendNotificationLogic) AdminSendNotification(ctx context.Context, req *types.AdminSendReq) (resp *types.AnyResp, err error) {
	adminID, _ := middleware.GetUserID(ctx)
	target := model.NotifyTargetUsers
	if req.SendAll {
		target = model.NotifyTargetAll
	}
	batch, err := biz.NewUserLogic(l.svcCtx).AdminSend(ctx, adminID, biz.AdminSendReq{
		Title:   req.Title,
		Content: req.Content,
		Target:  target,
		UserIDs: req.UserIDs,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: batch}, nil
}
