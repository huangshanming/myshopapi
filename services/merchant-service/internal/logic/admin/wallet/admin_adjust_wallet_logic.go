package wallet

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

type AdminAdjustWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAdjustWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAdjustWalletLogic {
	return &AdminAdjustWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminAdjustWalletLogic) AdminAdjustWallet(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	adminID, _ := middleware.GetUserID(ctx)
	var body types.WalletAdjustReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	wallet, err := biz.NewMerchantLogic(l.svcCtx).AdjustWallet(shopID, body.Field, body.Amount, body.Remark, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: wallet}, nil
}
