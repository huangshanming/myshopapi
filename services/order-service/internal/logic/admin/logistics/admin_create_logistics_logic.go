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

type AdminCreateLogisticsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateLogisticsLogic {
	return &AdminCreateLogisticsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminCreateLogisticsLogic) AdminCreateLogistics(ctx context.Context, req *types.LogisticsSaveBodyReq) (*types.LogisticsCompanyResp, error) {
	status := int8(req.Status)
	c, err := biz.NewLogisticsLogic(l.svcCtx).Create(ctx, types.LogisticsSaveReq{
		Name: req.Name, Code: req.Code, Sort: req.Sort, Status: &status,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.LogisticsCompanyResp{Data: c}, nil
}
