package favorite

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	huser "mymall/services/catalog-service/internal/product/app/user"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAddFavoriteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAddFavoriteLogic {
	return &UserAddFavoriteLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserAddFavoriteLogic) UserAddFavorite(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewFavoriteHandler(l.svcCtx).Add(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
