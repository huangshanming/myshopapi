package homepage

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

type AdminGrantSlotLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGrantSlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantSlotLogic {
	return &AdminGrantSlotLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantSlotLogic) AdminGrantSlot(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	adminID, _ := middleware.GetUserID(ctx)
	var body biz.GrantSlotReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	order, err := biz.NewMerchantLogic(l.svcCtx).GrantSlot(adminID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: order}, nil
}
