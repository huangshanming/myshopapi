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

type AdminDeleteArticleCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteArticleCategoryLogic(svcCtx *svc.ServiceContext) *AdminDeleteArticleCategoryLogic {
	return &AdminDeleteArticleCategoryLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteArticleCategoryLogic) AdminDeleteArticleCategory(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/admin/article-categories/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hadmin.NewArticleHandler(l.svcCtx).CategoryDelete)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
