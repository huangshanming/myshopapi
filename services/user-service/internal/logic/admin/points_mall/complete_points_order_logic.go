package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	data, err := hadmin.NewPointsOrderHandler(l.svcCtx).Complete(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
