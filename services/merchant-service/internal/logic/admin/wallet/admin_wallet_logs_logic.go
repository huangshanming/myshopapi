package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/merchant-service/internal/app/admin"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminWalletLogsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminWalletLogsLogic {
	return &AdminWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminWalletLogsLogic) AdminWalletLogs(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewWalletHandler(l.svcCtx).AdminWalletLogs(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.PageListResp
	if err := json.Unmarshal(b, &out); err != nil {
		var list interface{}
		if err2 := func() error { b,_:=json.Marshal(data); return json.Unmarshal(b, &list) }(); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
