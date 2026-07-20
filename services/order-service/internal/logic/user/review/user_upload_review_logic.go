package review

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type UserUploadReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUploadReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUploadReviewLogic {
	return &UserUploadReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUploadReviewLogic) UserUploadReview(w http.ResponseWriter, r *http.Request) {
	huser.NewReviewHandler(l.svcCtx).Upload(w, r)
}
