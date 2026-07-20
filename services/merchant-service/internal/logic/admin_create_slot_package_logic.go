package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateSlotPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCreateSlotPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateSlotPackageLogic {
	return &AdminCreateSlotPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateSlotPackageLogic) AdminCreateSlotPackage(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageSlotHandler(l.svcCtx).AdminCreateSlotPackage(w, r)
}
