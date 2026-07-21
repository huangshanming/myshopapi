package shop

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/upload"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUploadShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUploadShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUploadShopLogic {
	return &AdminUploadShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUploadShopLogic) AdminUploadShop(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Request: r}

	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := upload.SaveShopImage(shopID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: map[string]string{"url": url}}, nil
}
