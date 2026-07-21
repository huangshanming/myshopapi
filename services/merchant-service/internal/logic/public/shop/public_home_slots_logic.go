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

type PublicHomeSlotsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicHomeSlotsLogic(svcCtx *svc.ServiceContext) *PublicHomeSlotsLogic {
	return &PublicHomeSlotsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *PublicHomeSlotsLogic) PublicHomeSlots(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/shops/home-slots", nil, nil, nil, hpublic.NewHomepageSlotHandler(l.svcCtx).PublicHomeSlots)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
