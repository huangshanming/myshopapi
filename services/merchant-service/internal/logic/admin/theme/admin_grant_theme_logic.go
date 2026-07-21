package theme

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *AdminGrantThemeLogic) AdminGrantTheme(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	adminID, _ := middleware.GetUserID(ctx)
	var body biz.ThemeGrantReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	o, err := biz.NewMerchantLogic(l.svcCtx).GrantTheme(adminID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
