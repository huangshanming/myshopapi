package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"
	"net/url"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMineLogic {
	return &ListMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListMineLogic) ListMine(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	data, err := clogic.NewArticleLogic(l.svcCtx).UserListMine(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
