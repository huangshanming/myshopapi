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

type AdminPurgeArticleRecycleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminPurgeArticleRecycleLogic(svcCtx *svc.ServiceContext) *AdminPurgeArticleRecycleLogic {
	return &AdminPurgeArticleRecycleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminPurgeArticleRecycleLogic) AdminPurgeArticleRecycle(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/admin/articles/recycle", nil, nil, nil, hadmin.NewArticleHandler(l.svcCtx).RecycleDelete)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
