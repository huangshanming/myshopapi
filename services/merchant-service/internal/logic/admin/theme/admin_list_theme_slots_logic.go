package theme

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminListThemeSlotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListThemeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemeSlotsLogic {
	return &AdminListThemeSlotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemeSlotsLogic) AdminListThemeSlots(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageThemeHandler(l.svcCtx).AdminListThemeSlots(w, r)
}
