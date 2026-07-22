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

type AdminDeleteLogisticsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteLogisticsLogic {
	return &AdminDeleteLogisticsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminDeleteLogisticsLogic) AdminDeleteLogistics(ctx context.Context, req *types.IdPathReq) (*types.EmptyResp, error) {
	if err := biz.NewLogisticsLogic(l.svcCtx).Delete(ctx, req.Id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
