package points

import (
	"context"
	"encoding/json"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalRefundPointsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalRefundPointsLogic(svcCtx *svc.ServiceContext) *InternalRefundPointsLogic {
	return &InternalRefundPointsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalRefundPointsLogic) InternalRefundPoints(ctx context.Context, req *types.PointsLedgerReq) (resp *types.PointsResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/internal/points/refund", nil, nil, req, hinternal.NewTaskHandler(l.svcCtx).InternalRefundPoints)
	if err != nil {
		return nil, err
	}
	var out types.PointsResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		// may be bare number or {points:n}
		var n int64
		if err2 := json.Unmarshal(raw, &n); err2 == nil {
			out.Points = n
			return &out, nil
		}
		return nil, err
	}
	return &out, nil
}
