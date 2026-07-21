package article

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/catalog-service/internal/content/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminBatchAuditArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminBatchAuditArticlesLogic(svcCtx *svc.ServiceContext) *AdminBatchAuditArticlesLogic {
	return &AdminBatchAuditArticlesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminBatchAuditArticlesLogic) AdminBatchAuditArticles(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/articles/batch-audit", nil, nil, req, hadmin.NewArticleHandler(l.svcCtx).BatchAudit)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
