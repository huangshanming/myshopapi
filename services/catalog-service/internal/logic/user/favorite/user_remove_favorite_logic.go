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

func (l *UserRemoveFavoriteLogic) UserRemoveFavorite(ctx context.Context, req *types.ProductIdPathReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	if req.ProductId == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	if err := plogic.NewFavoriteLogic(l.svcCtx).Remove(ctx, userID, req.ProductId); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
