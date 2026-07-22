// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package points_mall

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadPointsProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadPointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPointsProductLogic {
	return &UploadPointsProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadPointsProductLogic) UploadPointsProduct() (resp *types.URLResp, err error) {
	// todo: add your logic here and delete this line

	return
}
