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

type CreateAdminLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(ctx context.Context, req *types.AdminCreateReq) (resp *types.AnyResp, err error) {
	user, err := biz.NewRBACLogic(l.svcCtx).CreateAdmin(ctx, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: user}, nil
}
