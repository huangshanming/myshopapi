package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/catalog-service/internal/content/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetArticleLogic {
	return &AdminGetArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetArticleLogic) AdminGetArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewArticleHandler(l.svcCtx).Detail(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
