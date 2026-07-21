package article

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hadmin "mymall/services/catalog-service/internal/content/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUploadArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUploadArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUploadArticleLogic {
	return &AdminUploadArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUploadArticleLogic) AdminUploadArticle(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewArticleHandler(l.svcCtx).Upload(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
