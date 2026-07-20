package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSlotSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSlotSettingsLogic {
	return &AdminListSlotSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotSettingsLogic) AdminListSlotSettings(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageSlotHandler(l.svcCtx).AdminListSlotSettings(w, r)
}
