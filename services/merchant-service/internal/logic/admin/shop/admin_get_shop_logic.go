package shop

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

type AdminGetShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetShopLogic {
	return &AdminGetShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetShopLogic) AdminGetShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	shop, err := biz.NewMerchantLogic(l.svcCtx).GetShop(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "店铺不存在")
	}
	return &types.AnyResp{Data: shop}, nil
}
