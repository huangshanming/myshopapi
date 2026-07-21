package comment

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

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

func (l *EmojiUpdateLogic) EmojiUpdate(ctx context.Context, req *types.EmojiUpdateBodyReq) (resp *types.AnyResp, err error) {
	sort := req.Sort
	if err := clogic.NewArticleLogic(l.svcCtx).UpdateEmoji(ctx, req.Id, req.Name, req.ImageURL, &sort, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
