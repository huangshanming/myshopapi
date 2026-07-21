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

type MerchantListThemeSlotsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListThemeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListThemeSlotsLogic {
	return &MerchantListThemeSlotsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListThemeSlotsLogic) MerchantListThemeSlots(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	list, err := biz.NewMerchantLogic(l.svcCtx).AdminListThemeSlots()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	on := make([]model.HomepageThemeSlot, 0)
	for _, s := range list {
		if s.Status == model.ThemeSlotOn {
			on = append(on, s)
		}
	}
	return &types.PageListResp{List: on}, nil
}
