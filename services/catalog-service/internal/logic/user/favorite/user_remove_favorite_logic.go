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

type UserRemoveFavoriteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserRemoveFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRemoveFavoriteLogic {
	return &UserRemoveFavoriteLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserRemoveFavoriteLogic) UserRemoveFavorite(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewFavoriteHandler(l.svcCtx).Remove(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
