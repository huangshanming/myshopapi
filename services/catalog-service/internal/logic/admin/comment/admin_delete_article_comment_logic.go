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

type AdminDeleteArticleCommentLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteArticleCommentLogic {
	return &AdminDeleteArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteArticleCommentLogic) AdminDeleteArticleComment(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).DeleteComment(ctx, id, 0); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
