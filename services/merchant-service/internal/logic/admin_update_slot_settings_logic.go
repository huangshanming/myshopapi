package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateSlotSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSlotSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateSlotSettingsLogic {
	return &AdminUpdateSlotSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSlotSettingsLogic) AdminUpdateSlotSettings(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageSlotHandler(l.svcCtx).AdminUpdateSlotSettings(w, r)
}
