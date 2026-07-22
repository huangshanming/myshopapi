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

type AdminCreateArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateArticleLogic {
	return &AdminCreateArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateArticleLogic) AdminCreateArticle(ctx context.Context, req *types.ArticleSaveReq) (resp *types.ArticleResp, err error) {
	uid, _ := middleware.GetUserID(ctx)
	a, err := clogic.NewArticleLogic(l.svcCtx).AdminCreate(ctx, uid, req.ToContent())
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.ArticleResp{Data: a}, nil
}
