package points

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

func ledgerFromReq(req *types.PointsLedgerReq, defaultChange string) biz.PointsLedgerReq {
	refID, _ := strconv.ParseUint(strings.TrimSpace(req.RefNo), 10, 64)
	changeType := strings.TrimSpace(req.ChangeType)
	if changeType == "" {
		changeType = defaultChange
	}
	refType := strings.TrimSpace(req.RefType)
	if refType == "" {
		refType = "ref"
	}
	return biz.PointsLedgerReq{
		UserID:     req.UserID,
		Points:     int(req.Points),
		ChangeType: changeType,
		Remark:     req.Reason,
		RefType:    refType,
		RefID:      refID,
	}
}

type InternalDeductPointsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalDeductPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalDeductPointsLogic {
	return &InternalDeductPointsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalDeductPointsLogic) InternalDeductPoints(ctx context.Context, req *types.PointsLedgerReq) (resp *types.PointsResp, err error) {
	p, err := biz.NewTaskLogic(l.svcCtx).DeductPoints(ctx, ledgerFromReq(req, model.PointChangeAdminAdjust))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsResp{Points: int64(p.Points)}, nil
}

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
	p, err := biz.NewTaskLogic(l.svcCtx).RefundPoints(ctx, ledgerFromReq(req, model.PointChangeAdminAdjust))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsResp{Points: int64(p.Points)}, nil
}

type InternalAddPointsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalAddPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalAddPointsLogic {
	return &InternalAddPointsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalAddPointsLogic) InternalAddPoints(ctx context.Context, req *types.PointsLedgerReq) (resp *types.PointsResp, err error) {
	// 加积分与 refund 同路径（正数入账），默认 change_type 区分业务
	p, err := biz.NewTaskLogic(l.svcCtx).RefundPoints(ctx, ledgerFromReq(req, model.PointChangeLotteryReward))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsResp{Points: int64(p.Points)}, nil
}
