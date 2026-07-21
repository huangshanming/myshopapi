package points_mall

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreatePointsProductLogic(svcCtx *svc.ServiceContext) *CreatePointsProductLogic {
	return &CreatePointsProductLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *CreatePointsProductLogic) CreatePointsProduct(ctx context.Context, req *types.PointsProductSaveReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/points-products", nil, nil, req, hadmin.NewPointsProductHandler(l.svcCtx).Create)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
