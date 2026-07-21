package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewJobStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobStatusLogic {
	return &JobStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *JobStatusLogic) JobStatus(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	job, err := plogic.NewProductAdminLogic(l.svcCtx).Job(ctx, shopID, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "任务不存在")
	}
	return &types.AnyResp{Data: job}, nil
}
