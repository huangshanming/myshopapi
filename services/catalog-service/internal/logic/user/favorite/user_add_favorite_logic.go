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

type UserAddFavoriteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserAddFavoriteLogic(svcCtx *svc.ServiceContext) *UserAddFavoriteLogic {
	return &UserAddFavoriteLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserAddFavoriteLogic) UserAddFavorite(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/favorites", nil, nil, req, huser.NewFavoriteHandler(l.svcCtx).Add)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
