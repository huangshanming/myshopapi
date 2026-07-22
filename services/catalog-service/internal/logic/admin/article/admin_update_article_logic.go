package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateArticleLogic {
	return &AdminUpdateArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateArticleLogic) AdminUpdateArticle(ctx context.Context, req *types.ArticleUpdateBodyReq) (resp *types.EmptyResp, err error) {
	uid, _ := middleware.GetUserID(ctx)
	if err := clogic.NewArticleLogic(l.svcCtx).AdminUpdate(ctx, req.Id, uid, req.ToContent()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
