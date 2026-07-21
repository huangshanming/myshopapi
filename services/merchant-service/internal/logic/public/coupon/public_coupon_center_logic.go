package coupon

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicCouponCenterLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicCouponCenterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicCouponCenterLogic {
	return &PublicCouponCenterLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicCouponCenterLogic) PublicCouponCenter(ctx context.Context) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{}

	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 && in.Request != nil {
		if raw := in.Request.Header.Get("X-User-Id"); raw != "" {
			userID, _ = strconv.ParseUint(raw, 10, 64)
		}
	}
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	list, err := biz.NewMerchantLogic(l.svcCtx).ListCenter(userID, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{"list": list}}, nil
}
