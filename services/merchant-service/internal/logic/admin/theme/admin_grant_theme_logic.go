package theme

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/merchant-service/internal/app/admin"
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
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewHomepageThemeHandler(l.svcCtx).AdminGrantTheme(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
