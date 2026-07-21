package theme

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

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

func (l *AdminUpdateThemeSlotLogic) AdminUpdateThemeSlot(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var body map[string]interface{}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	allowed := map[string]bool{
		"name": true, "desc": true, "cover_url": true, "default_link_type": true,
		"default_link_id": true, "status": true, "sort": true, "position": true,
	}
	updates := map[string]interface{}{}
	for k, v := range body {
		if allowed[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		return nil, xerr.New(http.StatusBadRequest, "无更新字段")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).AdminUpdateThemeSlot(id, updates); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
