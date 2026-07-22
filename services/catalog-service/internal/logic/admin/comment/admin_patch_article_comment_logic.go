package comment

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminPatchArticleCommentLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminPatchArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPatchArticleCommentLogic {
	return &AdminPatchArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminPatchArticleCommentLogic) AdminPatchArticleComment(ctx context.Context, req *types.ArticleCommentPatchBodyReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).PatchComment(ctx, id, 0, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
