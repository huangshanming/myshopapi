package article

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hmerchant "mymall/services/catalog-service/internal/content/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUploadArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUploadArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUploadArticleLogic {
	return &MerchantUploadArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUploadArticleLogic) MerchantUploadArticle(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hmerchant.NewArticleHandler(l.svcCtx).Upload(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
