package theme

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/model"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateThemePackageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateThemePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateThemePackageLogic {
	return &AdminCreateThemePackageLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateThemePackageLogic) AdminCreateThemePackage(ctx context.Context, req *types.ThemePackageSaveReq) (resp *types.AnyResp, err error) {
p := model.HomepageThemePackage{
		ThemeSlotID: req.ThemeSlotID, Name: req.Name, Price: req.Price, DurationDays: req.DurationDays,
		Status: req.Status, Sort: req.Sort, Remark: req.Remark,
	}
	if err := biz.NewMerchantLogic(l.svcCtx).AdminCreateThemePackage(&p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: p}, nil
}
