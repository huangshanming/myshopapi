package coupon

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

type ListMyCouponsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListMyCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyCouponsLogic {
	return &ListMyCouponsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListMyCouponsLogic) ListMyCoupons(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "请先登录")
	}
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListMyCoupons(userID, "" /* was query:status */, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
