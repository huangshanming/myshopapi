package address

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
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

func NewInternalGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalGetLogic {
	return &InternalGetLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalGetLogic) InternalGet(ctx context.Context, req *types.InternalAddressReq) (resp *types.AnyResp, err error) {
	data, err := hinternal.NewAddressHandler(l.svcCtx).InternalGet(ctx, appinput.CallInput{Query: url.Values{"id": {fmt.Sprintf("%d", req.Id)}, "user_id": {fmt.Sprintf("%d", req.UserID)}}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
