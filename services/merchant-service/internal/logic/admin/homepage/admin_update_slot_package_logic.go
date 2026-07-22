package homepage

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateSlotPackageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSlotPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateSlotPackageLogic {
	return &AdminUpdateSlotPackageLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSlotPackageLogic) AdminUpdateSlotPackage(ctx context.Context, req *types.SlotPackageUpdateBodyReq) (resp *types.EmptyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	p := model.HomepageSlotPackage{
		SlotType: req.SlotType, Name: req.Name, Price: req.Price, DurationDays: req.DurationDays,
		Status: req.Status, Sort: req.Sort, Remark: req.Remark,
	}
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateSlotPackage(req.Id, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
