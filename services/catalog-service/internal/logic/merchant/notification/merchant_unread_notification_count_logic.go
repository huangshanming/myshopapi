package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	nlogic "mymall/services/catalog-service/internal/notify/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUnreadNotificationCountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUnreadNotificationCountLogic {
	return &MerchantUnreadNotificationCountLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUnreadNotificationCountLogic) MerchantUnreadNotificationCount(ctx context.Context) (resp *types.CountResp, err error) {

	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	count, err := nlogic.NewNotificationLogic(l.svcCtx).UnreadCount(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.CountResp{Count: count}, nil
}
