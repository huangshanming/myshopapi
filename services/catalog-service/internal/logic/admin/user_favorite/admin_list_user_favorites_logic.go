package user_favorite

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListUserFavoritesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListUserFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListUserFavoritesLogic {
	return &AdminListUserFavoritesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListUserFavoritesLogic) AdminListUserFavorites(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	userID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || userID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	if pageSize <= 0 {
		pageSize = 50
	}
	list, total, err := plogic.NewFavoriteLogic(l.svcCtx).List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: map[string]interface{}{"list": list, "total": total}}, nil

}
