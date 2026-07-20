package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListApplicationsLogic {
	return &AdminListApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListApplicationsLogic) AdminListApplications(w http.ResponseWriter, r *http.Request) {
	hadmin.NewShopHandler(l.svcCtx).AdminListApplications(w, r)
}
