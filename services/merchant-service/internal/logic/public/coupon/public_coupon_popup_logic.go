package coupon

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicCouponPopupLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicCouponPopupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicCouponPopupLogic {
	return &PublicCouponPopupLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicCouponPopupLogic) PublicCouponPopup(ctx context.Context) (resp *types.AnyResp, err error) {

	userID, _ := middleware.GetUserID(ctx)
	list, err := biz.NewMerchantLogic(l.svcCtx).ListPopup(userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{"list": list}}, nil
}
