package theme

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListThemeSlotsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListThemeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemeSlotsLogic {
	return &AdminListThemeSlotsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemeSlotsLogic) AdminListThemeSlots(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	list, err := biz.NewMerchantLogic(l.svcCtx).AdminListThemeSlots()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
