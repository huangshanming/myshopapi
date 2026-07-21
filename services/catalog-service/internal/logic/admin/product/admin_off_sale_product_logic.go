package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	ptypes "mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOffSaleProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOffSaleProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOffSaleProductLogic {
	return &AdminOffSaleProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminOffSaleProductLogic) AdminOffSaleProduct(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	uid, _ := middleware.GetUserID(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body ptypes.PlatformProductRemarkReq
	_ = appinput.BindBody(in, &body)
	if err := plogic.NewPlatformProductLogic(l.svcCtx).ForceOffSale(ctx, id, uid, body.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
