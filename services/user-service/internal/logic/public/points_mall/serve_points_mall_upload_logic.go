package points_mall

import (
	"context"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServePointsMallUploadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewServePointsMallUploadLogic(svcCtx *svc.ServiceContext) *ServePointsMallUploadLogic {
	return &ServePointsMallUploadLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ServePointsMallUploadLogic) ServePointsMallUpload(ctx context.Context, req *types.FilePathReq) error {
	_ = req
	_ = ctx
	return nil
}
