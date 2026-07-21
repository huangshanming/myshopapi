package article

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hpublic "mymall/services/catalog-service/internal/content/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ListCommentsLogic) ListComments(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/articles/:id/comments", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hpublic.NewArticleHandler(l.svcCtx).ListComments)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		var list interface{}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
