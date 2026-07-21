package article

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	ctypes "mymall/services/catalog-service/internal/content/types"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMineLogic {
	return &CreateMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateMineLogic) CreateMine(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
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
	a, err := clogic.NewArticleLogic(l.svcCtx).UserCreate(ctx, userID, ctypes.ArticleSaveReq{
		CategoryID: body.CategoryID, Title: body.Title, CoverURL: body.CoverURL,
		Content: body.Content, ImageURLs: body.ImageURLs,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: a}, nil
}
