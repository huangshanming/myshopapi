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

type DeleteMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMineLogic {
	return &DeleteMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteMineLogic) DeleteMine(ctx context.Context, req *types.IdPathReq) (resp *types.EmptyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).UserDelete(ctx, userID, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
