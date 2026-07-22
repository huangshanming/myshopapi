package banner

import (
	"context"
	"io"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUploadBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadBannerLogic {
	return &UploadBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UploadBannerLogic) UploadBanner(ctx context.Context, r *http.Request) (resp *types.URLResp, err error) {
	if r == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	file, hdr, err := r.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := clogic.NewArticleLogic(l.svcCtx).SaveBannerUpload(hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.URLResp{Url: url}, nil
}
