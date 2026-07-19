package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	notifymodel "mymall/services/catalog-service/internal/notify/model"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/svc"
)

type PlatformProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlatformProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlatformProductLogic {
	return &PlatformProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlatformProductLogic) List(f repository.ProductListFilter) (map[string]interface{}, error) {
	f.PlatformScope = true
	list, total, err := l.svcCtx.ProductAdmin.List(f)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *PlatformProductLogic) ForceOffSale(id, operatorID uint64, remark string) error {
	remark = strings.TrimSpace(remark)
	if remark == "" {
		return errors.New("请填写备注，将通知商家")
	}
	p, err := l.svcCtx.ProductAdmin.GetByID(id)
	if err != nil {
		return errors.New("商品不存在")
	}
	if p.Status == model.ProductDeleted {
		return errors.New("商品已删除")
	}
	if err := l.svcCtx.ProductAdmin.SetStatus(id, p.ShopID, model.ProductOffSale); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]string{"status": model.ProductOffSale, "remark": remark})
	_ = l.svcCtx.ProductAdmin.AddOpLog(p.ShopID, &id, operatorID, "platform_off_sale", "", string(after))
	_ = l.svcCtx.Notifications.Create(&notifymodel.ShopNotification{
		ShopID:  p.ShopID,
		Type:    notifymodel.NotifProductOffSale,
		Title:   "商品被平台强制下架",
		Content: fmt.Sprintf("商品「%s」(ID:%d) 已被平台下架。备注：%s", p.Name, p.ID, remark),
		Link:    "/merchant/products",
		RefType: "product",
		RefID:   p.ID,
	})
	return nil
}

func (l *PlatformProductLogic) SoftDelete(id, operatorID uint64, remark string) error {
	remark = strings.TrimSpace(remark)
	if remark == "" {
		return errors.New("请填写备注，将通知商家")
	}
	p, err := l.svcCtx.ProductAdmin.GetByID(id)
	if err != nil {
		return errors.New("商品不存在")
	}
	if p.Status == model.ProductDeleted {
		return errors.New("商品已在回收站")
	}
	if err := l.svcCtx.ProductAdmin.SetStatus(id, p.ShopID, model.ProductDeleted); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]string{"status": model.ProductDeleted, "remark": remark})
	_ = l.svcCtx.ProductAdmin.AddOpLog(p.ShopID, &id, operatorID, "platform_delete", "", string(after))
	_ = l.svcCtx.Notifications.Create(&notifymodel.ShopNotification{
		ShopID:  p.ShopID,
		Type:    notifymodel.NotifProductDeleted,
		Title:   "商品被平台删除",
		Content: fmt.Sprintf("商品「%s」(ID:%d) 已被平台移入回收站。备注：%s", p.Name, p.ID, remark),
		Link:    "/merchant/products/recycle",
		RefType: "product",
		RefID:   p.ID,
	})
	return nil
}
