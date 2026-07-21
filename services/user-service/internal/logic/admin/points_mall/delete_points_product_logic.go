package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	_, err := hadmin.NewPointsProductHandler(l.svcCtx).Delete(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}})
	if err != nil {
		return err
	}
	return nil
}
