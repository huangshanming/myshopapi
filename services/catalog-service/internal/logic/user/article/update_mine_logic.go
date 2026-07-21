package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	ctypes "mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMineLogic {
	return &UpdateMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMineLogic) UpdateMine(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	var body struct {
		CategoryID uint64   `json:"category_id"`
		Title      string   `json:"title"`
		CoverURL   string   `json:"cover_url"`
		Content    string   `json:"content"`
		ImageURLs  []string `json:"image_urls"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).UserUpdate(ctx, userID, id, ctypes.ArticleSaveReq{
		CategoryID: body.CategoryID, Title: body.Title, CoverURL: body.CoverURL,
		Content: body.Content, ImageURLs: body.ImageURLs,
	}); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
