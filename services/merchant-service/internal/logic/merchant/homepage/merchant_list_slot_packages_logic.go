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

type MerchantListSlotPackagesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListSlotPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListSlotPackagesLogic {
	return &MerchantListSlotPackagesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListSlotPackagesLogic) MerchantListSlotPackages(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	list, err := biz.NewMerchantLogic(l.svcCtx).ListSlotPackages("" /* was query:slot_type */, true)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
