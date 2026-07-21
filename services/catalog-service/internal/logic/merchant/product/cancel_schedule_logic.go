package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelScheduleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCancelScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelScheduleLogic {
	return &CancelScheduleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CancelScheduleLogic) CancelSchedule(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewProductHandler(l.svcCtx).CancelSchedule(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
