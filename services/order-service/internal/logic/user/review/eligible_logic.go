package review

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EligibleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewEligibleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EligibleLogic {
	return &EligibleLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *EligibleLogic) Eligible(ctx context.Context, req *types.IdPathReq) (*types.ReviewEligibleResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	data, err := biz.NewReviewLogic(l.svcCtx).ReviewEligible(ctx, userID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return data, nil
}
