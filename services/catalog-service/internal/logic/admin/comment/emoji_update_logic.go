package comment

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

type EmojiUpdateLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewEmojiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiUpdateLogic {
	return &EmojiUpdateLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *EmojiUpdateLogic) EmojiUpdate(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
		Sort     *int   `json:"sort"`
		Status   *int8  `json:"status"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).UpdateEmoji(ctx, id, body.Name, body.ImageURL, body.Sort, body.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
