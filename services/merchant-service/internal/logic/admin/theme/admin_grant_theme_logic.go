package theme

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGrantThemeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGrantThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantThemeLogic {
	return &AdminGrantThemeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantThemeLogic) AdminGrantTheme(ctx context.Context, req *types.ThemeGrantReq) (resp *types.AnyResp, err error) {
	adminID, _ := middleware.GetUserID(ctx)
	o, err := biz.NewMerchantLogic(l.svcCtx).GrantTheme(adminID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
