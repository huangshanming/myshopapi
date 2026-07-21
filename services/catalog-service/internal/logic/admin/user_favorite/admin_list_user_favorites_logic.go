package user_favorite

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListUserFavoritesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListUserFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListUserFavoritesLogic {
	return &AdminListUserFavoritesLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListUserFavoritesLogic) AdminListUserFavorites(ctx context.Context, req *types.IdPageReq) (resp *types.PageListResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	list, total, err := plogic.NewFavoriteLogic(l.svcCtx).List(ctx, req.Id, req.Page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: map[string]interface{}{"list": list, "total": total}}, nil
}
