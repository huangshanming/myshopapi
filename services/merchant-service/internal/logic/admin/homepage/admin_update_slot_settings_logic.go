package homepage

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hadmin "mymall/services/merchant-service/internal/app/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateSlotSettingsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSlotSettingsLogic(svcCtx *svc.ServiceContext) *AdminUpdateSlotSettingsLogic {
	return &AdminUpdateSlotSettingsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSlotSettingsLogic) AdminUpdateSlotSettings(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/homepage-settings", nil, nil, req, hadmin.NewHomepageSlotHandler(l.svcCtx).AdminUpdateSlotSettings)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
