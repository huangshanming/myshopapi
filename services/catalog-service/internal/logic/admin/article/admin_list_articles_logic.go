package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticlesLogic {
	return &AdminListArticlesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListArticlesLogic) AdminListArticles(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	page, pageSize := in.Page()
	f := repository.ArticleListFilter{
		Title:       in.QueryGet("title"),
		AuditStatus: in.QueryGet("audit_status"),
		Status:      in.QueryGet("status"),
		Page:        page, PageSize: pageSize,
		Recycle: in.QueryGet("recycle") == "1" || (in.Request != nil && strings.Contains(in.Request.URL.Path, "/recycle")),
	}
	if s := in.QueryGet("shop_id"); s != "" {
		shopID, _ := strconv.ParseUint(s, 10, 64)
		f.ShopID = shopID
		f.FilterShop = true
	}
	if s := in.QueryGet("has_schedule"); s == "1" {
		v := true
		f.HasSchedule = &v
	} else if s == "0" {
		v := false
		f.HasSchedule = &v
	}
	if s := in.QueryGet("created_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := in.QueryGet("created_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	data, err := clogic.NewArticleLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
