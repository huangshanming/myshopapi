package application

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListApplicationsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListApplicationsLogic {
	return &AdminListApplicationsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListApplicationsLogic) AdminListApplications(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	p, ps := req.Page, req.PageSize
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListApplications(ctx, "" /* was query:status */, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
