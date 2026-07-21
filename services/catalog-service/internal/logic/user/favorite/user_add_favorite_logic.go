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

func (l *UserAddFavoriteLogic) UserAddFavorite(ctx context.Context, req *types.FavoriteAddReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	if err := plogic.NewFavoriteLogic(l.svcCtx).Add(ctx, userID, req.ProductID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
