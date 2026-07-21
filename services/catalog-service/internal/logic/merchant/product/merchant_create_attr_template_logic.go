package product

import (
	"context"
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

type MerchantCreateAttrTemplateLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateAttrTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateAttrTemplateLogic {
	return &MerchantCreateAttrTemplateLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateAttrTemplateLogic) MerchantCreateAttrTemplate(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var body ptypes.AttrTemplateReq
	_ = appinput.BindBody(in, &body)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	t := &model.ProductAttrTemplate{ID: id, ShopID: shopID, Name: body.Name, AttrsJSON: body.AttrsJSON, Status: 1}
	if err := l.svcCtx.ProductAdmin.SaveAttrTemplate(ctx, t); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: t}, nil
}
