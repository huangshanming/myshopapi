package points_mall

import (
	"context"
	"io"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadPointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUploadPointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPointsProductLogic {
	return &UploadPointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UploadPointsProductLogic) UploadPointsProduct(ctx context.Context) (resp *types.URLResp, err error) {
	return nil, xerr.New(400, "multipart required")
}

func (l *UploadPointsProductLogic) UploadFile(ctx context.Context, r io.Reader, filename string) (*types.URLResp, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	url, err := biz.NewPointsProductLogic(l.svcCtx).SaveUpload(ctx, filename, data)
	if err != nil {
		return nil, xerr.New(400, err.Error())
	}
	return &types.URLResp{Url: url}, nil
}
