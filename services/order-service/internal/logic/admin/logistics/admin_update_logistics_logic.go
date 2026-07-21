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

type AdminUpdateLogisticsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateLogisticsLogic {
	return &AdminUpdateLogisticsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminUpdateLogisticsLogic) AdminUpdateLogistics(ctx context.Context, req *types.LogisticsUpdateBodyReq) (*types.EmptyResp, error) {
	status := int8(req.Status)
	if err := biz.NewLogisticsLogic(l.svcCtx).Update(ctx, req.Id, types.LogisticsSaveReq{
		Name: req.Name, Code: req.Code, Sort: req.Sort, Status: &status,
	}); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
