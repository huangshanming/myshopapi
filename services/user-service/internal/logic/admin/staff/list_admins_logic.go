package staff

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ListAdminsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListAdminsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdminsLogic {
	return &ListAdminsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListAdminsLogic) ListAdmins(ctx context.Context, req *types.ListAdminsReq) (resp *types.PageListResp, err error) {
	list, total, err := biz.NewRBACLogic(l.svcCtx).ListAdmins(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
