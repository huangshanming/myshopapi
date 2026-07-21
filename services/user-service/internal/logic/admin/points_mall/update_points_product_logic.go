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

type UpdatePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePointsProductLogic {
	return &UpdatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdatePointsProductLogic) UpdatePointsProduct(ctx context.Context, req *types.PointsProductUpdateReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewPointsProductHandler(l.svcCtx).Update(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
