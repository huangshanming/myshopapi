package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	nlogic "mymall/services/catalog-service/internal/notify/logic"
	"mymall/services/catalog-service/internal/notify/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListNotificationsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListNotificationsLogic {
	return &MerchantListNotificationsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListNotificationsLogic) MerchantListNotifications(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := req.Page, req.PageSize
	f := repository.NotificationListFilter{ShopID: shopID, Page: page, PageSize: pageSize}
	if s := "" /* was query:is_read */; s == "0" || s == "1" {
		v := int8(0)
		if s == "1" {
			v = 1
		}
		f.IsRead = &v
	}
	data, err := nlogic.NewNotificationLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
