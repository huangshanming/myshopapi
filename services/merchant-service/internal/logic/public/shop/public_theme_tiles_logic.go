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

type PublicThemeTilesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicThemeTilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicThemeTilesLogic {
	return &PublicThemeTilesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicThemeTilesLogic) PublicThemeTiles(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewHomepageThemeHandler(l.svcCtx).PublicThemeTiles(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
