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

type UnlikeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUnlikeLogic(svcCtx *svc.ServiceContext) *UnlikeLogic {
	return &UnlikeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UnlikeLogic) Unlike(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/articles/:id/like", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hpublic.NewArticleHandler(l.svcCtx).Unlike)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
