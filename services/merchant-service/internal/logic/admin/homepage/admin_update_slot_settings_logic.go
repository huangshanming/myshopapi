package homepage

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/model"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateSlotSettingsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSlotSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateSlotSettingsLogic {
	return &AdminUpdateSlotSettingsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSlotSettingsLogic) AdminUpdateSlotSettings(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body struct {
		Items []model.HomepageSlotSetting `json:"items"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateSlotSettings(body.Items); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
