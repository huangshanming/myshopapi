// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompletePointsOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompletePointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompletePointsOrderLogic {
	return &CompletePointsOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompletePointsOrderLogic) CompletePointsOrder(req *types.IdPathReq) (resp *types.PointsOrderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
