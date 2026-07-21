package task

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserClaimLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserClaimLogic(svcCtx *svc.ServiceContext) *UserClaimLogic {
	return &UserClaimLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserClaimLogic) UserClaim(ctx context.Context, req *types.CodePathReq) (resp *types.PointsResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/tasks/{Code}/claim", map[string]string{"code": fmt.Sprintf("%v", req.Code)}, nil, nil, huser.NewTaskHandler(l.svcCtx).UserClaim)
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
