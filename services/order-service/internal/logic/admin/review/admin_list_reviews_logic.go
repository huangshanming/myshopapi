package review

import (
	"context"
	"net/http"

	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListReviewsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListReviewsLogic {
	return &AdminListReviewsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListReviewsLogic) AdminListReviews(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewReviewLogic(l.svcCtx).AdminList(ctx, 0, "", page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
