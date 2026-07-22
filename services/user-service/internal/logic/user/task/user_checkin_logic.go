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

type UserCheckinLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCheckinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCheckinLogic {
	return &UserCheckinLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserCheckinLogic) UserCheckin(ctx context.Context) (*types.CheckinResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	p, err := biz.NewTaskLogic(l.svcCtx).Checkin(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CheckinResp{Data: p}, nil
}
