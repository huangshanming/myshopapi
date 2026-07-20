package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type UploadPointsProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadPointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPointsProductLogic {
	return &UploadPointsProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadPointsProductLogic) UploadPointsProduct(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsProductHandler(l.svcCtx).Upload(w, r)
}
