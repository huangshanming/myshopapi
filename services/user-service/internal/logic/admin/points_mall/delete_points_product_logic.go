package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type DeletePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeletePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePointsProductLogic {
	return &DeletePointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeletePointsProductLogic) DeletePointsProduct(ctx context.Context, req *types.IdPathReq) error {
	if req.Id == 0 {
		return xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	if err := biz.NewPointsProductLogic(l.svcCtx).Delete(ctx, req.Id); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
