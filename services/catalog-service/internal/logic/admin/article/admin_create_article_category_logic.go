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

type AdminCreateArticleCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateArticleCategoryLogic(svcCtx *svc.ServiceContext) *AdminCreateArticleCategoryLogic {
	return &AdminCreateArticleCategoryLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateArticleCategoryLogic) AdminCreateArticleCategory(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/article-categories", nil, nil, req, hadmin.NewArticleHandler(l.svcCtx).CategoryCreate)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
