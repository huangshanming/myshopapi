package comment

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type EmojiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmojiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiListLogic {
	return &EmojiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmojiListLogic) EmojiList(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).EmojiList(w, r)
}
