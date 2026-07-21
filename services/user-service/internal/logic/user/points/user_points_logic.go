package points

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserPointsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPointsLogic {
	return &UserPointsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserPointsLogic) UserPoints(ctx context.Context) (*types.PointsResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	n, err := biz.NewTaskLogic(l.svcCtx).GetPoints(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PointsResp{Points: n.Points}, nil
}
