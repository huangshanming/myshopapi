package shop

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyShopsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMyShopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyShopsLogic {
	return &MyShopsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MyShopsLogic) MyShops(ctx context.Context) (resp *types.ShopListResp, err error) {

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	shops, err := biz.NewMerchantLogic(l.svcCtx).MyShops(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	out := make([]map[string]interface{}, 0, len(shops))
	for _, s := range shops {
		imgs, _ := l.svcCtx.Repo.ListShopImages(ctx, s.ID)
		urls := make([]string, 0, len(imgs))
		for _, im := range imgs {
			urls = append(urls, im.URL)
		}
		out = append(out, map[string]interface{}{
			"id": s.ID, "name": s.Name, "logo": s.Logo, "contact_name": s.ContactName,
			"contact_phone": s.ContactPhone, "description": s.Description, "category": s.Category,
			"province": s.Province, "city": s.City, "district": s.District, "address": s.Address,
			"latitude": s.Latitude, "longitude": s.Longitude, "local_enabled": s.LocalEnabled,
			"storefront_image": s.StorefrontImage, "status": s.Status, "images": urls,
		})
	}
	return &types.ShopListResp{Data: out}, nil
}
