package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminBatchAuditArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminBatchAuditArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminBatchAuditArticlesLogic {
	return &AdminBatchAuditArticlesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminBatchAuditArticlesLogic) AdminBatchAuditArticles(ctx context.Context, req *types.ArticleBatchAuditReq) (resp *types.AnyResp, err error) {
	if err := clogic.NewArticleLogic(l.svcCtx).BatchAudit(ctx, req.ToContent()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
