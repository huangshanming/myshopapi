// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPointsProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPointsProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointsProductsLogic {
	return &ListPointsProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPointsProductsLogic) ListPointsProducts(req *types.ListPointsProductsReq) (resp *types.PageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
