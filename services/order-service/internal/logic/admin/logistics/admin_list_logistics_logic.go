package logistics

import (
	"context"
	"net/http"

	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListLogisticsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListLogisticsLogic {
	return &AdminListLogisticsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListLogisticsLogic) AdminListLogistics(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewLogisticsLogic(l.svcCtx).List(ctx, repository.LogisticsListFilter{Page: page, PageSize: pageSize})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
