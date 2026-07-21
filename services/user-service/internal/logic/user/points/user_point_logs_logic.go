package points

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserPointLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserPointLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPointLogsLogic {
	return &UserPointLogsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserPointLogsLogic) UserPointLogs(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewTaskLogic(l.svcCtx).ListPointLogs(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
