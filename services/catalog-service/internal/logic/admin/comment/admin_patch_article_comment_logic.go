package comment

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

type AdminPatchArticleCommentLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminPatchArticleCommentLogic(svcCtx *svc.ServiceContext) *AdminPatchArticleCommentLogic {
	return &AdminPatchArticleCommentLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminPatchArticleCommentLogic) AdminPatchArticleComment(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PATCH", "/api/v1/admin/article-comments/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewArticleHandler(l.svcCtx).CommentPatch)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
