package wallet

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hadmin "mymall/services/merchant-service/internal/app/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminWalletLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminWalletLogsLogic(svcCtx *svc.ServiceContext) *AdminWalletLogsLogic {
	return &AdminWalletLogsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminWalletLogsLogic) AdminWalletLogs(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/shops/:id/wallet/logs", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hadmin.NewWalletHandler(l.svcCtx).AdminWalletLogs)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		var list interface{}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
