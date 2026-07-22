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

type AdminUpdateThemeSlotLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateThemeSlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateThemeSlotLogic {
	return &AdminUpdateThemeSlotLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateThemeSlotLogic) AdminUpdateThemeSlot(ctx context.Context, req *types.ThemeSlotUpdateBodyReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Desc != "" {
		updates["desc"] = req.Desc
	}
	if req.CoverURL != "" {
		updates["cover_url"] = req.CoverURL
	}
	if req.DefaultLinkType != "" {
		updates["default_link_type"] = req.DefaultLinkType
	}
	if req.DefaultLinkID != 0 {
		updates["default_link_id"] = req.DefaultLinkID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if len(updates) == 0 {
		return nil, xerr.New(http.StatusBadRequest, "无更新字段")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).AdminUpdateThemeSlot(id, updates); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
