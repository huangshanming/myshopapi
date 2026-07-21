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

type ListMyLikesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListMyLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLikesLogic {
	return &ListMyLikesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListMyLikesLogic) ListMyLikes(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, pageSize := int(req.Page), int(req.PageSize)
	data, err := clogic.NewArticleLogic(l.svcCtx).ListMyLikes(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
