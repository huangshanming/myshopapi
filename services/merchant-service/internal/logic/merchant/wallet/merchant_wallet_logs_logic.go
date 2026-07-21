package wallet

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"net/url"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantWalletLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantWalletLogsLogic {
	return &MerchantWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantWalletLogsLogic) MerchantWalletLogs(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	shopID := middleware.GetShopID(ctx)
	p, ps := in.Page()
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListWalletLogs(shopID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
