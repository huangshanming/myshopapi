package coupon

import (
	"context"
	"fmt"
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

type MerchantCopyCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCopyCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCopyCouponLogic {
	return &MerchantCopyCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCopyCouponLogic) MerchantCopyCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	c, err := biz.NewMerchantLogic(l.svcCtx).CopyCoupon(id, shopID, userID, false)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: c}, nil
}
