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

type GetByOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetByOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByOrderLogic {
	return &GetByOrderLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *GetByOrderLogic) GetByOrder(ctx context.Context, req *types.IdPathReq) (*types.AnyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	rev, err := biz.NewReviewLogic(l.svcCtx).GetByOrder(ctx, userID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.AnyResp{Data: rev}, nil
}
