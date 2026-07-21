package logistics

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/order-service/internal/app/admin"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateLogisticsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateLogisticsLogic(svcCtx *svc.ServiceContext) *AdminCreateLogisticsLogic {
	return &AdminCreateLogisticsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateLogisticsLogic) AdminCreateLogistics(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/logistics", nil, nil, req, hadmin.NewLogisticsHandler(l.svcCtx).Create)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
