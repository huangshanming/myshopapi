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

type PublicGetShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicGetShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicGetShopLogic {
	return &PublicGetShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicGetShopLogic) PublicGetShop(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	shop, err := biz.NewMerchantLogic(l.svcCtx).GetPublicShop(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: shop}, nil
}
