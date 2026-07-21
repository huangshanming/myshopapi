package notification

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	nlogic "mymall/services/catalog-service/internal/notify/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantMarkNotificationReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantMarkNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantMarkNotificationReadLogic {
	return &MerchantMarkNotificationReadLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantMarkNotificationReadLogic) MerchantMarkNotificationRead(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := nlogic.NewNotificationLogic(l.svcCtx).MarkRead(ctx, id, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
