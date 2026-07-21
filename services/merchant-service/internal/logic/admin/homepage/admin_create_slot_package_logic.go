package homepage

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/model"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateSlotPackageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateSlotPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateSlotPackageLogic {
	return &AdminCreateSlotPackageLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateSlotPackageLogic) AdminCreateSlotPackage(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var p model.HomepageSlotPackage
	if err := appinput.BindBody(in, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).CreateSlotPackage(&p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: p}, nil
}
