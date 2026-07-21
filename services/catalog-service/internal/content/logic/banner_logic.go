package logic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/content/model"
	"mymall/services/catalog-service/internal/uploadpath"
)

type BannerSaveReq struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkType string `json:"link_type"`
	LinkID   uint64 `json:"link_id"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
}

func parseBannerTime(s string) (*common.LocalTime, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			lt := common.LocalTime(t)
			return &lt, nil
		}
	}
	return nil, errors.New("时间格式无效")
}

func (l *ArticleLogic) normalizeBanner(req BannerSaveReq) (*model.HomepageBanner, error) {
	img := strings.TrimSpace(req.ImageURL)
	if img == "" {
		return nil, errors.New("请上传 Banner 图片")
	}
	linkType := strings.TrimSpace(req.LinkType)
	if linkType == "" {
		linkType = model.BannerLinkNone
	}
	switch linkType {
	case model.BannerLinkNone:
		req.LinkID = 0
	case model.BannerLinkProduct:
		if req.LinkID == 0 || !l.svcCtx.Articles.ProductExistsOnSale(context.Background(), req.LinkID) {
			return nil, errors.New("请选择有效在售商品")
		}
	case model.BannerLinkArticle:
		if req.LinkID == 0 || !l.svcCtx.Articles.ArticleExistsPublished(context.Background(), req.LinkID) {
			return nil, errors.New("请选择已发布文章")
		}
	default:
		return nil, errors.New("跳转类型无效")
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = model.BannerOn
	}
	if status != model.BannerOn && status != model.BannerOff {
		return nil, errors.New("状态无效")
	}
	startAt, err := parseBannerTime(req.StartAt)
	if err != nil {
		return nil, err
	}
	endAt, err := parseBannerTime(req.EndAt)
	if err != nil {
		return nil, err
	}
	return &model.HomepageBanner{
		Title:    strings.TrimSpace(req.Title),
		ImageURL: img,
		LinkType: linkType,
		LinkID:   req.LinkID,
		Sort:     req.Sort,
		Status:   status,
		StartAt:  startAt,
		EndAt:    endAt,
	}, nil
}

func (l *ArticleLogic) PublicBanners() ([]model.HomepageBanner, error) {
	return l.svcCtx.Articles.ListBannersPublic(context.Background())
}

func (l *ArticleLogic) AdminListBanners(page, pageSize int) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Articles.ListBannersAdmin(context.Background(), page, pageSize)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Articles.FillBannerLinkNames(context.Background(), list)
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *ArticleLogic) AdminGetBanner(id uint64) (*model.HomepageBanner, error) {
	b, err := l.svcCtx.Articles.GetBanner(context.Background(), id)
	if err != nil {
		return nil, errors.New("Banner 不存在")
	}
	tmp := []model.HomepageBanner{*b}
	l.svcCtx.Articles.FillBannerLinkNames(context.Background(), tmp)
	b.LinkName = tmp[0].LinkName
	return b, nil
}

func (l *ArticleLogic) AdminCreateBanner(req BannerSaveReq) (*model.HomepageBanner, error) {
	b, err := l.normalizeBanner(req)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Articles.CreateBanner(context.Background(), b); err != nil {
		return nil, err
	}
	return b, nil
}

func (l *ArticleLogic) AdminUpdateBanner(id uint64, req BannerSaveReq) error {
	if _, err := l.svcCtx.Articles.GetBanner(context.Background(), id); err != nil {
		return errors.New("Banner 不存在")
	}
	b, err := l.normalizeBanner(req)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"title":     b.Title,
		"image_url": b.ImageURL,
		"link_type": b.LinkType,
		"link_id":   b.LinkID,
		"sort":      b.Sort,
		"status":    b.Status,
		"start_at":  b.StartAt,
		"end_at":    b.EndAt,
	}
	return l.svcCtx.Articles.UpdateBanner(context.Background(), id, updates)
}

func (l *ArticleLogic) AdminDeleteBanner(id uint64) error {
	return l.svcCtx.Articles.DeleteBanner(context.Background(), id)
}

func (l *ArticleLogic) SaveBannerUpload(filename string, data []byte) (string, error) {
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	dir := uploadpath.Abs("banners")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/banners/" + name, nil
}
