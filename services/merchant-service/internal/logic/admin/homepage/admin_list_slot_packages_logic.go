package homepage

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

type AdminListSlotPackagesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotPackagesLogic(svcCtx *svc.ServiceContext) *AdminListSlotPackagesLogic {
	return &AdminListSlotPackagesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotPackagesLogic) AdminListSlotPackages(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/homepage-packages", nil, url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}, nil, hadmin.NewHomepageSlotHandler(l.svcCtx).AdminListSlotPackages)
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
