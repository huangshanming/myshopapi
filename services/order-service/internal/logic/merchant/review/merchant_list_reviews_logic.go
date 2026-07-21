package review

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListReviewsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListReviewsLogic {
	return &MerchantListReviewsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantListReviewsLogic) MerchantListReviews(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	shopID := middleware.GetShopID(ctx)
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewReviewLogic(l.svcCtx).MerchantList(ctx, shopID, "", page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
