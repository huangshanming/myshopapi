package favorite

import (
	"context"
	"fmt"
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

type UserFavoriteStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserFavoriteStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFavoriteStatusLogic {
	return &UserFavoriteStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserFavoriteStatusLogic) UserFavoriteStatus(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	productID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || productID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	okFav, err := plogic.NewFavoriteLogic(l.svcCtx).IsFavorited(ctx, userID, productID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]bool{"favorited": okFav}}, nil
}
