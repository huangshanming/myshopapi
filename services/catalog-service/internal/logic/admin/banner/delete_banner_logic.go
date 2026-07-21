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

type DeleteBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBannerLogic {
	return &DeleteBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteBannerLogic) DeleteBanner(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).AdminDeleteBanner(id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
