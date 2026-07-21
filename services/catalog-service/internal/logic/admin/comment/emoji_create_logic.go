package comment

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmojiCreateLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewEmojiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiCreateLogic {
	return &EmojiCreateLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *EmojiCreateLogic) EmojiCreate(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
		Sort     int    `json:"sort"`
		Status   *int8  `json:"status"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	status := int8(1)
	if body.Status != nil {
		status = *body.Status
	}
	e, err := clogic.NewArticleLogic(l.svcCtx).CreateEmoji(ctx, body.Name, body.ImageURL, body.Sort, status)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: e}, nil
}
