package homepage

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminUpdateSlotPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSlotPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateSlotPackageLogic {
	return &AdminUpdateSlotPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSlotPackageLogic) AdminUpdateSlotPackage(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageSlotHandler(l.svcCtx).AdminUpdateSlotPackage(w, r)
}
