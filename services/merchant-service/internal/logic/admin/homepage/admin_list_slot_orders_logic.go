package homepage

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"net/url"
	"strconv"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSlotOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSlotOrdersLogic {
	return &AdminListSlotOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotOrdersLogic) AdminListSlotOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	p, ps := in.Page()
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListSlotOrders(shopID, in.QueryGet("slot_type"), in.QueryGet("status"), p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
