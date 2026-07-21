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

func (l *CancelPointsOrderLogic) CancelPointsOrder(ctx context.Context, req *types.RemarkReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewPointsOrderHandler(l.svcCtx).Cancel(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
