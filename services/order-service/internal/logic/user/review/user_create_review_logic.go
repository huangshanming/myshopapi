package review

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type UserCreateReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateReviewLogic {
	return &UserCreateReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateReviewLogic) UserCreateReview(w http.ResponseWriter, r *http.Request) {
	huser.NewReviewHandler(l.svcCtx).Create(w, r)
}
