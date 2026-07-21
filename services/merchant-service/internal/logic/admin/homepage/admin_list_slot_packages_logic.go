package homepage

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSlotPackagesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSlotPackagesLogic {
	return &AdminListSlotPackagesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotPackagesLogic) AdminListSlotPackages(ctx context.Context, req *types.SlotTypePageReq) (resp *types.PageListResp, err error) {
	list, err := biz.NewMerchantLogic(l.svcCtx).ListSlotPackages(req.SlotType, false)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
