package article

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	ctypes "mymall/services/catalog-service/internal/content/types"
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

func (l *AdminBatchAuditArticlesLogic) AdminBatchAuditArticles(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body ctypes.ArticleBatchAuditReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).BatchAudit(ctx, body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
