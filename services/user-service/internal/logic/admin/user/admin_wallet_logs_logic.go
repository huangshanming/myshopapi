package user

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/url"

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

func (l *AdminWalletLogsLogic) AdminWalletLogs(ctx context.Context, req *types.AdminWalletLogsReq) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/users/{Id}/wallet/logs", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}, nil, hadmin.NewWalletHandler(l.svcCtx).AdminWalletLogs)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		// raw may already be {list,total}
		var m map[string]json.RawMessage
		if err2 := json.Unmarshal(raw, &m); err2 == nil {
			_ = json.Unmarshal(m["list"], &out.List)
			_ = json.Unmarshal(m["total"], &out.Total)
			return &out, nil
		}
		return nil, err
	}
	return &out, nil
}
