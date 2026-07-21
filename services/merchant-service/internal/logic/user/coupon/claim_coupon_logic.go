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

type ClaimCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewClaimCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ClaimCouponLogic) ClaimCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "请先登录")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body struct {
		Source string `json:"source"`
	}
	_ = appinput.BindBody(in, &body)
	uc, err := biz.NewMerchantLogic(l.svcCtx).ClaimCoupon(userID, id, body.Source)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: uc}, nil
}
