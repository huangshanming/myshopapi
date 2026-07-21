package user_favorite

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/catalog-service/internal/product/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListUserFavoritesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListUserFavoritesLogic(svcCtx *svc.ServiceContext) *AdminListUserFavoritesLogic {
	return &AdminListUserFavoritesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminListUserFavoritesLogic) AdminListUserFavorites(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/users/:id/favorites", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hadmin.NewFavoriteHandler(l.svcCtx).AdminUserList)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		var list interface{}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
