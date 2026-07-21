package shop

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/merchant-service/internal/app/public"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicHomeSlotsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicHomeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicHomeSlotsLogic {
	return &PublicHomeSlotsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicHomeSlotsLogic) PublicHomeSlots(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewHomepageSlotHandler(l.svcCtx).PublicHomeSlots(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
