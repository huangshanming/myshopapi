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

type PublicThemeTilesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicThemeTilesLogic(svcCtx *svc.ServiceContext) *PublicThemeTilesLogic {
	return &PublicThemeTilesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *PublicThemeTilesLogic) PublicThemeTiles(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/home/theme-tiles", nil, nil, nil, hpublic.NewHomepageThemeHandler(l.svcCtx).PublicThemeTiles)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
