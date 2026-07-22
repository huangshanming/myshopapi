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

type AdminDeleteProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteProductLogic {
	return &AdminDeleteProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteProductLogic) AdminDeleteProduct(ctx context.Context, req *types.PlatformProductRemarkBodyReq) (resp *types.EmptyResp, err error) {
	uid, _ := middleware.GetUserID(ctx)
	id := req.Id
	if err := plogic.NewPlatformProductLogic(l.svcCtx).SoftDelete(ctx, id, uid, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
