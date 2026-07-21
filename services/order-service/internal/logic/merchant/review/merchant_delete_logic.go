package review

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hmerchant "mymall/services/order-service/internal/app/merchant"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantDeleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantDeleteLogic(svcCtx *svc.ServiceContext) *MerchantDeleteLogic {
	return &MerchantDeleteLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantDeleteLogic) MerchantDelete(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/merchant/reviews/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hmerchant.NewReviewHandler(l.svcCtx).MerchantDelete)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
