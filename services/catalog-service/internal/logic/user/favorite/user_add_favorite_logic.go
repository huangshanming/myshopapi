package favorite

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *UserAddFavoriteLogic) UserAddFavorite(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	var body struct {
		ProductID uint64 `json:"product_id"`
	}
	if err := appinput.BindBody(in, &body); err != nil || body.ProductID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := plogic.NewFavoriteLogic(l.svcCtx).Add(ctx, userID, body.ProductID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
