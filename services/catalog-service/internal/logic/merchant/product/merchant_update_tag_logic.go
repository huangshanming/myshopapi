package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	ptypes "mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUpdateTagLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUpdateTagLogic {
	return &MerchantUpdateTagLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUpdateTagLogic) MerchantUpdateTag(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var body ptypes.TagReq
	_ = appinput.BindBody(in, &body)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	tag := &model.ProductTag{ID: id, ShopID: shopID, Name: body.Name, Color: body.Color, Status: 1}
	if err := l.svcCtx.ProductAdmin.SaveTag(ctx, tag); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: tag}, nil
}
