package homepage

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/merchant-service/internal/app/admin"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGrantSlotLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGrantSlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantSlotLogic {
	return &AdminGrantSlotLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantSlotLogic) AdminGrantSlot(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewHomepageSlotHandler(l.svcCtx).AdminGrantSlot(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
