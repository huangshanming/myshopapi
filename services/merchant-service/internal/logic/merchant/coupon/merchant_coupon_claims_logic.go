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

type MerchantCouponClaimsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponClaimsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponClaimsLogic {
	return &MerchantCouponClaimsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponClaimsLogic) MerchantCouponClaims(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := biz.NewMerchantLogic(l.svcCtx).CouponClaims(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{"list": list, "total": total}}, nil

}
