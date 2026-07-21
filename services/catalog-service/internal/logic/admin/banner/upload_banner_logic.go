package banner

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hadmin "mymall/services/catalog-service/internal/content/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUploadBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadBannerLogic {
	return &UploadBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UploadBannerLogic) UploadBanner(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewArticleHandler(l.svcCtx).UploadBanner(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
