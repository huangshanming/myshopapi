package order

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetOrderLogic {
	return &UserGetOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserGetOrderLogic) UserGetOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewOrderHandler(l.svcCtx).Detail(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
