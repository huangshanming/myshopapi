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

type DetailPointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDetailPointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailPointsProductLogic {
	return &DetailPointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DetailPointsProductLogic) DetailPointsProduct(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewPointsProductHandler(l.svcCtx).Detail(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
