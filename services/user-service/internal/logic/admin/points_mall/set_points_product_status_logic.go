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

type SetPointsProductStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSetPointsProductStatusLogic(svcCtx *svc.ServiceContext) *SetPointsProductStatusLogic {
	return &SetPointsProductStatusLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *SetPointsProductStatusLogic) SetPointsProductStatus(ctx context.Context, req *types.PointsProductStatusReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/points-products/{Id}/status", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewPointsProductHandler(l.svcCtx).SetStatus)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
