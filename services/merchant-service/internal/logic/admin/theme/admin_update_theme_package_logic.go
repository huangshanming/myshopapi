package theme

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateThemePackageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateThemePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateThemePackageLogic {
	return &AdminUpdateThemePackageLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateThemePackageLogic) AdminUpdateThemePackage(ctx context.Context, req *types.ThemePackageUpdateBodyReq) (resp *types.AnyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	updates := map[string]interface{}{
		"theme_slot_id": req.ThemeSlotID, "name": req.Name, "price": req.Price,
		"duration_days": req.DurationDays, "status": req.Status, "sort": req.Sort, "remark": req.Remark,
	}
	if err := biz.NewMerchantLogic(l.svcCtx).AdminUpdateThemePackage(req.Id, updates); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
