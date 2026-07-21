package product

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAttrTemplatesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListAttrTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAttrTemplatesLogic {
	return &ListAttrTemplatesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListAttrTemplatesLogic) ListAttrTemplates(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewProductHandler(l.svcCtx).ListAttrTemplates(ctx, appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}})
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
