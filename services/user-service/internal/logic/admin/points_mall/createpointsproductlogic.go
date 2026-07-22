// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePointsProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePointsProductLogic {
	return &CreatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePointsProductLogic) CreatePointsProduct(req *types.PointsProductSaveReq) (resp *types.PointsProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
