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

type AdminAuditArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAuditArticleLogic(svcCtx *svc.ServiceContext) *AdminAuditArticleLogic {
	return &AdminAuditArticleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminAuditArticleLogic) AdminAuditArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/articles/:id/audit", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewArticleHandler(l.svcCtx).Audit)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
