package task

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserClaimLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClaimLogic {
	return &UserClaimLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserClaimLogic) UserClaim(ctx context.Context, req *types.CodePathReq) (resp *types.PointsResp, err error) {
	data, err := huser.NewTaskHandler(l.svcCtx).UserClaim(ctx, appinput.CallInput{PathVars: map[string]string{"code": fmt.Sprintf("%v", req.Code)}})
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
