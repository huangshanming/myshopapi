package coupon

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalMatchCouponsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalMatchCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalMatchCouponsLogic {
	return &InternalMatchCouponsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalMatchCouponsLogic) InternalMatchCoupons(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body biz.MatchReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := biz.NewMerchantLogic(l.svcCtx).MatchCoupons(body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
