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

type LikeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *LikeLogic) Like(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/articles/:id/like", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hpublic.NewArticleHandler(l.svcCtx).Like)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
