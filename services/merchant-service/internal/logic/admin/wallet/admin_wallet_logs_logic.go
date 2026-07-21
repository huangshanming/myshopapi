package wallet

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

type AdminWalletLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminWalletLogsLogic {
	return &AdminWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminWalletLogsLogic) AdminWalletLogs(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	p, ps := in.Page()
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListWalletLogs(shopID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
