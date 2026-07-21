package theme

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

type AdminListThemePackagesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListThemePackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemePackagesLogic {
	return &AdminListThemePackagesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemePackagesLogic) AdminListThemePackages(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	slotID, _ := strconv.ParseUint(in.QueryGet("theme_slot_id"), 10, 64)
	list, err := biz.NewMerchantLogic(l.svcCtx).ListThemePackages(slotID, false)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
