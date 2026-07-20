package biz

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/uploadpath"

	"gorm.io/gorm"
)

type ReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewLogic {
	return &ReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReviewLogic) ReviewEligible(userID, orderID uint64) (map[string]interface{}, error) {
	order, err := l.svcCtx.Repo.FindByID(orderID, userID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	exists, _ := l.svcCtx.Reviews.ExistsByOrderID(orderID)
	return map[string]interface{}{
		"eligible": order.Status == model.OrderStatusCompleted && !exists,
		"reviewed": exists || order.Status == model.OrderStatusReviewed,
		"status":   order.Status,
		"order_id": order.ID,
	}, nil
}

func (l *ReviewLogic) Create(userID, orderID uint64, req model.CreateReviewReq) (*model.ProductReview, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("评分须为 1-5")
	}
	if utf8.RuneCountInString(req.Content) > 1000 {
		return nil, errors.New("评价内容过长")
	}
	if len(req.Images) > 9 {
		return nil, errors.New("图片最多 9 张")
	}
	order, err := l.svcCtx.Repo.FindByID(orderID, userID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order.Status != model.OrderStatusCompleted {
		return nil, errors.New("仅确认收货后的订单可评价")
	}
	exists, err := l.svcCtx.Reviews.ExistsByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该订单已评价")
	}
	items := order.Items
	if len(items) == 0 {
		return nil, errors.New("订单商品不存在")
	}
	var item *model.OrderItem
	if req.OrderItemID > 0 {
		for i := range items {
			if items[i].ID == req.OrderItemID {
				item = &items[i]
				break
			}
		}
		if item == nil {
			return nil, errors.New("订单项不属于该订单")
		}
	} else {
		item = &items[0]
	}
	rev := &model.ProductReview{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		UserID:      userID,
		ShopID:      order.ShopID,
		ProductID:   item.ProductID,
		OrderItemID: item.ID,
		SkuID:       item.SkuID,
		SkuSnapshot: item.SkuSnapshot,
		Rating:      req.Rating,
		Content:     strings.TrimSpace(req.Content),
		IsAnonymous: req.IsAnonymous,
		Status:      model.ReviewStatusVisible,
	}
	if err := l.svcCtx.Reviews.Create(rev, req.Images); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单状态已变更，请刷新后重试")
		}
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "uk_order") {
			return nil, errors.New("该订单已评价")
		}
		return nil, err
	}
	_ = l.refreshProductStats(item.ProductID)
	full, _ := l.svcCtx.Reviews.GetByOrderID(orderID)
	if full != nil {
		return full, nil
	}
	return rev, nil
}

func (l *ReviewLogic) refreshProductStats(productID uint64) error {
	avg, count, goodRate, err := l.svcCtx.Reviews.ProductStats(productID)
	if err != nil {
		return err
	}
	return l.svcCtx.Reviews.UpdateProductStats(productID, avg, count, goodRate)
}

func (l *ReviewLogic) GetByOrder(userID, orderID uint64) (*model.ProductReview, error) {
	if _, err := l.svcCtx.Repo.FindByID(orderID, userID); err != nil {
		return nil, errors.New("订单不存在")
	}
	rev, err := l.svcCtx.Reviews.GetByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("暂无评价")
		}
		return nil, err
	}
	return rev, nil
}

func (l *ReviewLogic) ListByProduct(productID uint64, page, pageSize int) ([]model.ProductReview, int64, error) {
	list, total, err := l.svcCtx.Reviews.ListByProduct(productID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].IsAnonymous {
			list[i].UserName = "匿名用户"
			list[i].UserID = 0
		}
	}
	return list, total, nil
}

func (l *ReviewLogic) MerchantList(shopID uint64, ratingLevel string, page, pageSize int) ([]model.ProductReview, int64, error) {
	return l.svcCtx.Reviews.ListMerchant(shopID, ratingLevel, page, pageSize)
}

func (l *ReviewLogic) AdminList(shopID uint64, ratingLevel string, page, pageSize int) ([]model.ProductReview, int64, error) {
	return l.svcCtx.Reviews.ListAdmin(shopID, ratingLevel, page, pageSize)
}

func (l *ReviewLogic) Reply(shopID, reviewID uint64, reply string) error {
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return errors.New("回复不能为空")
	}
	if utf8.RuneCountInString(reply) > 500 {
		return errors.New("回复过长")
	}
	if err := l.svcCtx.Reviews.Reply(reviewID, shopID, reply); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评价不存在")
		}
		return err
	}
	return nil
}

func (l *ReviewLogic) SoftDelete(reviewID uint64, shopID uint64) error {
	rev, err := l.svcCtx.Reviews.GetByID(reviewID)
	if err != nil {
		return errors.New("评价不存在")
	}
	if err := l.svcCtx.Reviews.SoftDelete(reviewID, shopID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评价不存在")
		}
		return err
	}
	_ = l.refreshProductStats(rev.ProductID)
	return nil
}

func (l *ReviewLogic) SaveUpload(userID uint64, filename string, data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("空文件")
	}
	if len(data) > 5*1024*1024 {
		return "", errors.New("图片不能超过 5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
	default:
		return "", errors.New("仅支持 jpg/png/webp/gif")
	}
	dir := uploadpath.Abs("reviews", fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("/uploads/reviews/%d/%s", userID, name), nil
}
