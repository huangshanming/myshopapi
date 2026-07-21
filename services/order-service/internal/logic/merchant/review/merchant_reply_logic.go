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

type MerchantReplyLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantReplyLogic {
	return &MerchantReplyLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantReplyLogic) MerchantReply(ctx context.Context, req *types.ReviewReplyBodyReq) (*types.EmptyResp, error) {
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewReviewLogic(l.svcCtx).Reply(ctx, shopID, req.Id, req.Reply); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
