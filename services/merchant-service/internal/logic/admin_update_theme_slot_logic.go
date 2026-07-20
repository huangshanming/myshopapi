package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateThemeSlotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateThemeSlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateThemeSlotLogic {
	return &AdminUpdateThemeSlotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateThemeSlotLogic) AdminUpdateThemeSlot(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageThemeHandler(l.svcCtx).AdminUpdateThemeSlot(w, r)
}
