package review

import (
	"context"
	"io"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUploadReviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserUploadReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUploadReviewLogic {
	return &UserUploadReviewLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserUploadReviewLogic) UserUploadReview(ctx context.Context, r *http.Request) (*types.URLResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "请上传文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取文件失败")
	}
	url, err := biz.NewReviewLogic(l.svcCtx).SaveUpload(userID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.URLResp{Url: url}, nil
}
