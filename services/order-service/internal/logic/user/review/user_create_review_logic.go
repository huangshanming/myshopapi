package review

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateReviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateReviewLogic {
	return &UserCreateReviewLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserCreateReviewLogic) UserCreateReview(ctx context.Context, req *types.CreateReviewBodyReq) (*types.AnyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	rev, err := biz.NewReviewLogic(l.svcCtx).Create(ctx, userID, req.Id, model.CreateReviewReq{
		Rating: req.Rating, Content: req.Content, IsAnonymous: req.IsAnonymous,
		OrderItemID: req.OrderItemID, Images: req.Images,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: rev}, nil
}
