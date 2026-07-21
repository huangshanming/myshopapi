package theme

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListThemePackagesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListThemePackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemePackagesLogic {
	return &AdminListThemePackagesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemePackagesLogic) AdminListThemePackages(ctx context.Context, req *types.ThemePackageListReq) (resp *types.PageListResp, err error) {
	slotID := req.ThemeSlotId
	list, err := biz.NewMerchantLogic(l.svcCtx).ListThemePackages(slotID, false)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
