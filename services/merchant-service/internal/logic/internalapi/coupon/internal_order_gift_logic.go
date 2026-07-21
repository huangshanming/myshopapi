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

type InternalOrderGiftLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalOrderGiftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalOrderGiftLogic {
	return &InternalOrderGiftLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalOrderGiftLogic) InternalOrderGift(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body struct {
		UserID uint64 `json:"user_id"`
		ShopID uint64 `json:"shop_id"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	n, err := biz.NewMerchantLogic(l.svcCtx).OrderGiftCoupons(body.UserID, body.ShopID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{"granted": n}}, nil
}
