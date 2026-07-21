package address

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalGetLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalGetLogic(svcCtx *svc.ServiceContext) *InternalGetLogic {
	return &InternalGetLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalGetLogic) InternalGet(ctx context.Context, req *types.InternalAddressReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/addresses/internal", nil, url.Values{"id": {fmt.Sprintf("%d", req.Id)}, "user_id": {fmt.Sprintf("%d", req.UserID)}}, nil, hinternal.NewAddressHandler(l.svcCtx).InternalGet)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
