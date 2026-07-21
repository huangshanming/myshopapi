package homepage

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"net/url"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSlotPackagesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSlotPackagesLogic {
	return &AdminListSlotPackagesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotPackagesLogic) AdminListSlotPackages(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	list, err := biz.NewMerchantLogic(l.svcCtx).ListSlotPackages(in.QueryGet("slot_type"), false)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
