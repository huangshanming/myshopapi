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

type AdminOfflineArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOfflineArticleLogic(svcCtx *svc.ServiceContext) *AdminOfflineArticleLogic {
	return &AdminOfflineArticleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminOfflineArticleLogic) AdminOfflineArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/articles/:id/offline", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewArticleHandler(l.svcCtx).Offline)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
