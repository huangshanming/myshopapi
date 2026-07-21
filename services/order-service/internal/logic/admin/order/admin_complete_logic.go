package order

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCompleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCompleteLogic {
	return &AdminCompleteLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminCompleteLogic) AdminComplete(ctx context.Context, req *types.IdPathReq) (*types.EmptyResp, error) {
	if err := biz.NewOrderLogic(l.svcCtx).Complete(ctx, req.Id, 0); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
