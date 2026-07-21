package order

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListOrdersLogic {
	return &UserListOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserListOrdersLogic) UserListOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewOrderHandler(l.svcCtx).List(ctx, appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.PageListResp
	if err := json.Unmarshal(b, &out); err != nil {
		var list interface{}
		if err2 := func() error { b, _ := json.Marshal(data); return json.Unmarshal(b, &list) }(); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
