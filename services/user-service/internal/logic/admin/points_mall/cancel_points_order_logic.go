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

type CancelPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCancelPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPointsOrderLogic {
	return &CancelPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CancelPointsOrderLogic) CancelPointsOrder(ctx context.Context, req *types.RemarkReq) (resp *types.PointsOrderResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	o, err := biz.NewPointsOrderLogic(l.svcCtx).AdminCancel(ctx, req.Id, req.Remark)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsOrderResp{Data: o}, nil
}
