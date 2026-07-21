package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailUserPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDetailUserPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailUserPointsOrderLogic {
	return &DetailUserPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DetailUserPointsOrderLogic) DetailUserPointsOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	data, err := huser.NewPointsOrderHandler(l.svcCtx).Detail(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
