package points

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type InternalRefundPointsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalRefundPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalRefundPointsLogic {
	return &InternalRefundPointsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalRefundPointsLogic) InternalRefundPoints(ctx context.Context, req *types.PointsLedgerReq) (resp *types.PointsResp, err error) {
	refID, _ := strconv.ParseUint(req.RefNo, 10, 64)
	p, err := biz.NewTaskLogic(l.svcCtx).RefundPoints(ctx, biz.PointsLedgerReq{
		UserID:     req.UserID,
		Points:     int(req.Points),
		ChangeType: model.PointChangeAdminAdjust,
		Remark:     req.Reason,
		RefType:    "ref",
		RefID:      refID,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsResp{Points: int64(p.Points)}, nil
}
