package comment

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type EmojiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmojiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiUpdateLogic {
	return &EmojiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmojiUpdateLogic) EmojiUpdate(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).EmojiUpdate(w, r)
}
