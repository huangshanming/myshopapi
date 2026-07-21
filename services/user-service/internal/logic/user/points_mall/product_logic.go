package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
)

type ProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductLogic {
	return &ProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductLogic) ListPointsProduct(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	list, total, err := l.svcCtx.PointsProducts.List(l.ctx, page, pageSize, model.PointsProductStatusOn, "")
	if err != nil {
		httpx.ErrorCtx(l.ctx, w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(l.ctx, w, map[string]interface{}{"list": list, "total": total})
}
