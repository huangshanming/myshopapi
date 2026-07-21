package health

import (
	"context"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadyzLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *ReadyzLogic) Readyz(ctx context.Context) (resp *types.EmptyResp, err error) {
	_ = ctx

	return &types.EmptyResp{}, nil
}
