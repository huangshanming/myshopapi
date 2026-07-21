package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailUserPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDetailUserPointsOrderLogic(svcCtx *svc.ServiceContext) *DetailUserPointsOrderLogic {
	return &DetailUserPointsOrderLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *DetailUserPointsOrderLogic) DetailUserPointsOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/points-mall/orders/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, huser.NewPointsOrderHandler(l.svcCtx).Detail)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
