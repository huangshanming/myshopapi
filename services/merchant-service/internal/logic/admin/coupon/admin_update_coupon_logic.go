package coupon

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateCouponLogic {
	return &AdminUpdateCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateCouponLogic) AdminUpdateCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body biz.CouponSaveReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateCoupon(id, 0, true, body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
