package favorite

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"
	"strconv"

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
	in := appinput.CallInput{}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	productID, err := strconv.ParseUint(in.Path("product_id"), 10, 64)
	if err != nil || productID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	if err := plogic.NewFavoriteLogic(l.svcCtx).Remove(ctx, userID, productID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
