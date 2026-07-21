package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type CompletePointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCompletePointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompletePointsOrderLogic {
	return &CompletePointsOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CompletePointsOrderLogic) CompletePointsOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	o, err := biz.NewPointsOrderLogic(l.svcCtx).AdminComplete(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
