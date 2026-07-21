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

type DeleteMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteMineLogic(svcCtx *svc.ServiceContext) *DeleteMineLogic {
	return &DeleteMineLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *DeleteMineLogic) DeleteMine(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/user/articles/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hpublic.NewArticleHandler(l.svcCtx).DeleteMine)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
