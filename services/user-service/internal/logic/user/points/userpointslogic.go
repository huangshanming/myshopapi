// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserPointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPointsLogic {
	return &UserPointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPointsLogic) UserPoints() (resp *types.PointsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
