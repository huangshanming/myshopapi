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

type RemarkPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRemarkPointsOrderLogic(svcCtx *svc.ServiceContext) *RemarkPointsOrderLogic {
	return &RemarkPointsOrderLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *RemarkPointsOrderLogic) RemarkPointsOrder(ctx context.Context, req *types.RemarkReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/points-orders/{Id}/remark", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewPointsOrderHandler(l.svcCtx).Remark)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
