package shop

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicHomeSlotsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicHomeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicHomeSlotsLogic {
	return &PublicHomeSlotsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicHomeSlotsLogic) PublicHomeSlots(ctx context.Context, req *types.SlotTypeQueryReq) (resp *types.AnyResp, err error) {
	slotType := req.SlotType
	list, err := biz.NewMerchantLogic(l.svcCtx).HomeSlots(slotType)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: list}, nil
}
