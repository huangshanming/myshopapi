package health

import (
	"context"
	"fmt"
	"net/http/httptest"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadyzLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ReadyzLogic) Readyz(ctx context.Context) (resp *types.EmptyResp, err error) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/readyz", nil).WithContext(ctx)
	l.svcCtx.Health.ReadyHandler().ServeHTTP(rr, req)
	if rr.Code >= 400 {
		return nil, fmt.Errorf("not ready: %s", rr.Body.String())
	}
	return &types.EmptyResp{}, nil
}
