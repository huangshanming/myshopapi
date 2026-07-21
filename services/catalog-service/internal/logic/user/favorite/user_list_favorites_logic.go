package favorite

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListFavoritesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListFavoritesLogic {
	return &UserListFavoritesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserListFavoritesLogic) UserListFavorites(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := plogic.NewFavoriteLogic(l.svcCtx).List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: map[string]interface{}{"list": list, "total": total}}, nil

}
