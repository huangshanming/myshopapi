package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDetailPointsOrderLogic(svcCtx *svc.ServiceContext) *DetailPointsOrderLogic {
	return &DetailPointsOrderLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *DetailPointsOrderLogic) DetailPointsOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/points-orders/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewPointsOrderHandler(l.svcCtx).Detail)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
