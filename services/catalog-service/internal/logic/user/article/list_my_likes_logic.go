package article

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/catalog-service/internal/content/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyLikesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListMyLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLikesLogic {
	return &ListMyLikesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListMyLikesLogic) ListMyLikes(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewArticleHandler(l.svcCtx).ListMyLikes(ctx, appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.PageListResp
	if err := json.Unmarshal(b, &out); err != nil {
		var list interface{}
		if err2 := func() error { b, _ := json.Marshal(data); return json.Unmarshal(b, &list) }(); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
