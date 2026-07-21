package shop

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

type AdminUpdateShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateShopLogic {
	return &AdminUpdateShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateShopLogic) AdminUpdateShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	var body types.AdminUpdateShopReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).AdminUpdateShop(ctx, id, body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
