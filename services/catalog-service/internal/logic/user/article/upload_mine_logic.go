package article

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hpublic "mymall/services/catalog-service/internal/content/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUploadMineLogic(svcCtx *svc.ServiceContext) *UploadMineLogic {
	return &UploadMineLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UploadMineLogic) UploadMine(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/article-uploads", nil, nil, req, hpublic.NewArticleHandler(l.svcCtx).UploadMine)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
