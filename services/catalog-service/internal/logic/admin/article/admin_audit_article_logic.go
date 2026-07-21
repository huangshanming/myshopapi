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

type AdminAuditArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAuditArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAuditArticleLogic {
	return &AdminAuditArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminAuditArticleLogic) AdminAuditArticle(ctx context.Context, req *types.ArticleAuditBodyReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).Audit(ctx, id, req.ToContent()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
