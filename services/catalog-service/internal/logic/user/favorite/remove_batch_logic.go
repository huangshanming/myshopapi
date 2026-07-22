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

type RemoveBatchLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRemoveBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBatchLogic {
	return &RemoveBatchLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *RemoveBatchLogic) RemoveBatch(ctx context.Context, req *types.FavoriteBatchRemoveReq) (resp *types.EmptyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	if err := plogic.NewFavoriteLogic(l.svcCtx).RemoveBatch(ctx, userID, req.ProductIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
