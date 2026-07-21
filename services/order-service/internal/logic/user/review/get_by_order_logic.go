package review

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetByOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetByOrderLogic(svcCtx *svc.ServiceContext) *GetByOrderLogic {
	return &GetByOrderLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *GetByOrderLogic) GetByOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/orders/:id/review", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, huser.NewReviewHandler(l.svcCtx).GetByOrder)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
