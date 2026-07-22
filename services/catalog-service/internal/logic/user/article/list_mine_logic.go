package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

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
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, pageSize := int(req.Page), int(req.PageSize)
	data, err := clogic.NewArticleLogic(l.svcCtx).UserListMine(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.FromPaged(data), nil
}
