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

type UserCreateOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateOrderLogic {
	return &UserCreateOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserCreateOrderLogic) UserCreateOrder(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewOrderHandler(l.svcCtx).List(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
