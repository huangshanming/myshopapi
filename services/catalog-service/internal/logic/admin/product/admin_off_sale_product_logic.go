package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOffSaleProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOffSaleProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOffSaleProductLogic {
	return &AdminOffSaleProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminOffSaleProductLogic) AdminOffSaleProduct(ctx context.Context, req *types.PlatformProductRemarkBodyReq) (resp *types.EmptyResp, err error) {
	uid, _ := middleware.GetUserID(ctx)
	id := req.Id
	if err := plogic.NewPlatformProductLogic(l.svcCtx).ForceOffSale(ctx, id, uid, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
