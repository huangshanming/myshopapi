package theme

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/model"
	"net/http"
	"strconv"

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

func (l *AdminUpdateThemePackageLogic) AdminUpdateThemePackage(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var p model.HomepageThemePackage
	if err := appinput.BindBody(in, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	updates := map[string]interface{}{
		"theme_slot_id": p.ThemeSlotID, "name": p.Name, "price": p.Price,
		"duration_days": p.DurationDays, "status": p.Status, "sort": p.Sort, "remark": p.Remark,
	}
	if err := biz.NewMerchantLogic(l.svcCtx).AdminUpdateThemePackage(id, updates); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
