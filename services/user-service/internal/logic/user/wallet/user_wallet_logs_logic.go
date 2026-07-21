package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserWalletLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserWalletLogsLogic(svcCtx *svc.ServiceContext) *UserWalletLogsLogic {
	return &UserWalletLogsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserWalletLogsLogic) UserWalletLogs(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/wallet/logs", nil, url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}, nil, huser.NewWalletHandler(l.svcCtx).UserWalletLogs)
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
