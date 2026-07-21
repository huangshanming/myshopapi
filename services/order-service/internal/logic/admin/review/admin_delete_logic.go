package review

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteLogic {
	return &AdminDeleteLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminDeleteLogic) AdminDelete(ctx context.Context, req *types.IdPathReq) (*types.AnyResp, error) {
	if err := biz.NewReviewLogic(l.svcCtx).SoftDelete(ctx, req.Id, 0); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
