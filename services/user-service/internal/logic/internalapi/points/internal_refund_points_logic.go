package points

import (
	"context"
	"encoding/json"
	"mymall/pkg/appinput"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	data, err := hinternal.NewTaskHandler(l.svcCtx).InternalRefundPoints(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.PointsResp
	if err := json.Unmarshal(b, &out); err != nil {
		// may be bare number or {points:n}
		var n int64
		if err2 := json.Unmarshal(b, &n); err2 == nil {
			out.Points = n
			return &out, nil
		}
		return nil, err
	}
	return &out, nil
}
