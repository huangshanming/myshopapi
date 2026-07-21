package points_mall

import (
	"context"
	"encoding/json"
	"fmt"
	"mymall/pkg/appinput"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPointsOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListPointsOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointsOrdersLogic {
	return &ListPointsOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListPointsOrdersLogic) ListPointsOrders(ctx context.Context, req *types.ListPointsOrdersReq) (resp *types.PageListResp, err error) {
	data, err := hadmin.NewPointsOrderHandler(l.svcCtx).List(ctx, appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}, "status": {req.Status}, "order_no": {req.OrderNo}, "keyword": {req.Keyword}, "user_id": {fmt.Sprintf("%d", req.UserID)}}})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.PageListResp
	if err := json.Unmarshal(b, &out); err != nil {
		// raw may already be {list,total}
		var m map[string]json.RawMessage
		if err2 := json.Unmarshal(b, &m); err2 == nil {
			_ = json.Unmarshal(m["list"], &out.List)
			_ = json.Unmarshal(m["total"], &out.Total)
			return &out, nil
		}
		return nil, err
	}
	return &out, nil
}
