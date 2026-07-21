package logistics

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UpdateStatusLogic) UpdateStatus(ctx context.Context, req *types.LogisticsStatusBodyReq) (*types.EmptyResp, error) {
	if err := biz.NewLogisticsLogic(l.svcCtx).UpdateStatus(ctx, req.Id, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
