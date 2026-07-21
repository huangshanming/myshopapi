package article

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hpublic "mymall/services/catalog-service/internal/content/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUploadMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMineLogic {
	return &UploadMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UploadMineLogic) UploadMine(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hpublic.NewArticleHandler(l.svcCtx).UploadMine(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
