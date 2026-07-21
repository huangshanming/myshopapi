package theme

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListThemeOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListThemeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemeOrdersLogic {
	return &AdminListThemeOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemeOrdersLogic) AdminListThemeOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopID, _ := strconv.ParseUint("" /* was query:shop_id */, 10, 64)
	slotID, _ := strconv.ParseUint("" /* was query:theme_slot_id */, 10, 64)
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListThemeOrders(shopID, slotID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
