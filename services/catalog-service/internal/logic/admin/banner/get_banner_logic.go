package banner

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBannerLogic {
	return &GetBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetBannerLogic) GetBanner(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	b, err := clogic.NewArticleLogic(l.svcCtx).AdminGetBanner(id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.AnyResp{Data: b}, nil
}
