package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	ctypes "mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOfflineArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOfflineArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOfflineArticleLogic {
	return &AdminOfflineArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminOfflineArticleLogic) AdminOfflineArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body ctypes.ArticleRemarkReq
	_ = appinput.BindBody(in, &body)
	if err := clogic.NewArticleLogic(l.svcCtx).Offline(ctx, id, body.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
