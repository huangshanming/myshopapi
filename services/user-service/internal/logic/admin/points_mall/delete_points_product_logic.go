package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeletePointsProductLogic(svcCtx *svc.ServiceContext) *DeletePointsProductLogic {
	return &DeletePointsProductLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *DeletePointsProductLogic) DeletePointsProduct(ctx context.Context, req *types.IdPathReq) error {
	_, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/admin/points-products/{Id}", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, nil, hadmin.NewPointsProductHandler(l.svcCtx).Delete)
	if err != nil {
		return err
	}
	return nil
}
