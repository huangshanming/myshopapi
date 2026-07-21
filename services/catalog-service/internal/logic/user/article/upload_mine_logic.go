package article

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUploadMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMineLogic {
	return &UploadMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UploadMineLogic) UploadMine(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Request: r}

	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	_ = userID
	if err := in.Request.ParseMultipartForm(6 << 20); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "上传失败")
	}
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "请选择文件")
	}
	defer file.Close()
	buf, err := io.ReadAll(io.LimitReader(file, 5<<20+1))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取文件失败")
	}
	url, err := clogic.NewArticleLogic(l.svcCtx).SaveUpload(0, hdr.Filename, buf)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: map[string]string{"url": url}}, nil
}
