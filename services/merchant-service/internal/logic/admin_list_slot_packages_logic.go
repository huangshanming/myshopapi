package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSlotPackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSlotPackagesLogic {
	return &AdminListSlotPackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotPackagesLogic) AdminListSlotPackages(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageSlotHandler(l.svcCtx).AdminListSlotPackages(w, r)
}
