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

type AdminUpdateSlotPackageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSlotPackageLogic(svcCtx *svc.ServiceContext) *AdminUpdateSlotPackageLogic {
	return &AdminUpdateSlotPackageLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSlotPackageLogic) AdminUpdateSlotPackage(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/homepage-packages/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewHomepageSlotHandler(l.svcCtx).AdminUpdateSlotPackage)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
