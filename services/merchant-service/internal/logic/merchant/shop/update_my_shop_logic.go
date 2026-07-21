package shop

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

type UpdateMyShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMyShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyShopLogic {
	return &UpdateMyShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyShopLogic) UpdateMyShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	var body types.UpdateShopReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateMyShop(ctx, shopID, userID, body); err != nil {
		return nil, xerr.New(http.StatusForbidden, err.Error())
	}
	return &types.AnyResp{}, nil
}
