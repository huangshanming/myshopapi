package shop

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hpublic "mymall/services/merchant-service/internal/app/public"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicGetShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicGetShopLogic(svcCtx *svc.ServiceContext) *PublicGetShopLogic {
	return &PublicGetShopLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *PublicGetShopLogic) PublicGetShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/shops/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hpublic.NewShopHandler(l.svcCtx).PublicGetShop)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
