package points_mall

import (
	"context"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePointsProductLogic {
	return &CreatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreatePointsProductLogic) CreatePointsProduct(ctx context.Context, req *types.PointsProductSaveReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewPointsProductHandler(l.svcCtx).Create(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
