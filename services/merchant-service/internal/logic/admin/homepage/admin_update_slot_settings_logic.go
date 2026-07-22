package homepage

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

func (l *AdminUpdateSlotSettingsLogic) AdminUpdateSlotSettings(ctx context.Context, req *types.UpdateSlotSettingsReq) (resp *types.EmptyResp, err error) {
items := make([]model.HomepageSlotSetting, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, model.HomepageSlotSetting{SlotType: it.SlotType, HomeLimit: it.HomeLimit})
	}
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateSlotSettings(items); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
