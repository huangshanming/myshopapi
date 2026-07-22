package homepage

import (
	"context"
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

func (l *AdminCreateSlotPackageLogic) AdminCreateSlotPackage(ctx context.Context, req *types.SlotPackageSaveReq) (resp *types.SlotPackageResp, err error) {
p := model.HomepageSlotPackage{
		SlotType: req.SlotType, Name: req.Name, Price: req.Price, DurationDays: req.DurationDays,
		Status: req.Status, Sort: req.Sort, Remark: req.Remark,
	}
	if err := biz.NewMerchantLogic(l.svcCtx).CreateSlotPackage(&p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.SlotPackageResp{Data: p}, nil
}
