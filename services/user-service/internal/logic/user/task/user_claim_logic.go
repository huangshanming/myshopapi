package task

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

type UserClaimLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClaimLogic {
	return &UserClaimLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserClaimLogic) UserClaim(ctx context.Context, req *types.CodePathReq) (*types.PointsResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	n, err := biz.NewTaskLogic(l.svcCtx).Claim(ctx, userID, req.Code)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsResp{Points: int64(n.Points)}, nil
}
