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

type UpdateMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMineLogic(svcCtx *svc.ServiceContext) *UpdateMineLogic {
	return &UpdateMineLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMineLogic) UpdateMine(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/user/articles/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hpublic.NewArticleHandler(l.svcCtx).UpdateMine)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
