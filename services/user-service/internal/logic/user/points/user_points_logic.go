package points

import (
	"context"
	"encoding/json"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserPointsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserPointsLogic(svcCtx *svc.ServiceContext) *UserPointsLogic {
	return &UserPointsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserPointsLogic) UserPoints(ctx context.Context) (resp *types.PointsResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/points", nil, nil, nil, huser.NewTaskHandler(l.svcCtx).UserPoints)
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
