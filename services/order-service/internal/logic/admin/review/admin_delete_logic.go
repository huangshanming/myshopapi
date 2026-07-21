package review

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/order-service/internal/app/admin"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogic(svcCtx *svc.ServiceContext) *AdminDeleteLogic {
	return &AdminDeleteLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteLogic) AdminDelete(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/admin/reviews/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hadmin.NewReviewHandler(l.svcCtx).AdminDelete)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
