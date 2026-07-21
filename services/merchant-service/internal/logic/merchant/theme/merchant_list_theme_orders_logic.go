package theme

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"net/url"
	"strconv"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListThemeOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListThemeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListThemeOrdersLogic {
	return &MerchantListThemeOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListThemeOrdersLogic) MerchantListThemeOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListThemeOrders(shopID, 0, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
