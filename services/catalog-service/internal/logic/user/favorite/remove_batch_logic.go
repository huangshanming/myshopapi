package favorite

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	huser "mymall/services/catalog-service/internal/product/app/user"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveBatchLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRemoveBatchLogic(svcCtx *svc.ServiceContext) *RemoveBatchLogic {
	return &RemoveBatchLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *RemoveBatchLogic) RemoveBatch(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/favorites/batch-remove", nil, nil, req, huser.NewFavoriteHandler(l.svcCtx).RemoveBatch)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
