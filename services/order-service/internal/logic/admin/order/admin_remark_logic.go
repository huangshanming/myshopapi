package order

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRemarkLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRemarkLogic {
	return &AdminRemarkLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminRemarkLogic) AdminRemark(ctx context.Context, req *types.RemarkBodyReq) (*types.EmptyResp, error) {
	if err := biz.NewOrderLogic(l.svcCtx).UpdateRemark(ctx, req.Id, 0, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
