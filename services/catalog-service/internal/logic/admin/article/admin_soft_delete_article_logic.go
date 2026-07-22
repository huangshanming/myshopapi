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

type AdminSoftDeleteArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminSoftDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminSoftDeleteArticleLogic {
	return &AdminSoftDeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminSoftDeleteArticleLogic) AdminSoftDeleteArticle(ctx context.Context, req *types.ArticleRemarkBodyReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).SoftDelete(ctx, id, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
