package types

import (
	"encoding/json"

	ctypes "mymall/services/catalog-service/internal/content/types"
	ptypes "mymall/services/catalog-service/internal/product/types"
	sotypes "mymall/services/catalog-service/internal/shopops/types"
)

func (r *MerchantProductSaveReq) ToProduct() ptypes.MerchantProductSaveReq {
	var out ptypes.MerchantProductSaveReq
	b, _ := json.Marshal(r)
	_ = json.Unmarshal(b, &out)
	return out
}

func (r *BatchProductReq) ToProduct() ptypes.BatchProductReq {
	var out ptypes.BatchProductReq
	b, _ := json.Marshal(r)
	_ = json.Unmarshal(b, &out)
	return out
}

func (r *BatchStockReq) ToProduct() ptypes.BatchStockReq {
	var out ptypes.BatchStockReq
	b, _ := json.Marshal(r)
	_ = json.Unmarshal(b, &out)
	return out
}

func (r *RecycleReq) ToProduct() ptypes.RecycleReq {
	return ptypes.RecycleReq{ProductIDs: r.ProductIDs}
}

func (r *TagReq) ToProduct() ptypes.TagReq {
	return ptypes.TagReq{Name: r.Name, Color: r.Color}
}

func (r *AttrTemplateReq) ToProduct() ptypes.AttrTemplateReq {
	return ptypes.AttrTemplateReq{Name: r.Name, AttrsJSON: r.AttrsJSON}
}

func (r *CategoryReq) ToProduct() ptypes.CategoryReq {
	var out ptypes.CategoryReq
	b, _ := json.Marshal(r)
	_ = json.Unmarshal(b, &out)
	return out
}

func (r *ArticleSaveReq) ToContent() ctypes.ArticleSaveReq {
	var out ctypes.ArticleSaveReq
	b, _ := json.Marshal(r)
	_ = json.Unmarshal(b, &out)
	return out
}

func (r *ArticleBatchAuditReq) ToContent() ctypes.ArticleBatchAuditReq {
	return ctypes.ArticleBatchAuditReq{IDs: r.IDs, Pass: r.Pass, RejectReason: r.RejectReason}
}

func (r *ArticleCategorySaveReq) ToContent() ctypes.ArticleCategorySaveReq {
	return ctypes.ArticleCategorySaveReq{ParentID: r.ParentID, Name: r.Name, Sort: r.Sort, Status: r.Status}
}

func (r *UserArticleCreateReq) ToContent() ctypes.ArticleSaveReq {
	return ctypes.ArticleSaveReq{
		CategoryID: r.CategoryID, Title: r.Title, CoverURL: r.CoverURL,
		Content: r.Content, ImageURLs: r.ImageURLs,
	}
}

func (r *ShopRoleReq) ToShopOps() sotypes.ShopRoleReq {
	return sotypes.ShopRoleReq{Code: r.Code, Name: r.Name, Remark: r.Remark, MenuIDs: r.MenuIDs}
}

func (r *ShopStaffReq) ToShopOps() sotypes.ShopStaffReq {
	return sotypes.ShopStaffReq{
		Mobile: r.Mobile, RoleID: r.RoleID, Nickname: r.Nickname, Password: r.Password, Mode: r.Mode,
	}
}

func (r *ProductUpdateBodyReq) ToProduct() ptypes.MerchantProductSaveReq {
	return r.MerchantProductSaveReq.ToProduct()
}

func (r *TagUpdateBodyReq) ToProduct() ptypes.TagReq {
	return ptypes.TagReq{Name: r.Name, Color: r.Color}
}

func (r *AttrTemplateUpdateBodyReq) ToProduct() ptypes.AttrTemplateReq {
	return ptypes.AttrTemplateReq{Name: r.Name, AttrsJSON: r.AttrsJSON}
}

func (r *StockAdjustBodyReq) ToProduct() ptypes.StockAdjustReq {
	return ptypes.StockAdjustReq{SkuID: r.SkuID, Stock: r.Stock, Delta: r.Delta}
}

func (r *ScheduleBodyReq) ToProduct() ptypes.ScheduleReq {
	return ptypes.ScheduleReq{Action: r.Action, RunAt: r.RunAt}
}

func (r *SetStatusBodyReq) ToProduct() ptypes.SetStatusReq {
	return ptypes.SetStatusReq{Status: r.Status}
}

func (r *CategoryUpdateBodyReq) ToProduct() ptypes.CategoryReq {
	tmp := CategoryReq{
		ParentId: r.ParentId, Name: r.Name, Icon: r.Icon, Description: r.Description,
		SortOrder: r.SortOrder, Level: r.Level, IsShow: r.IsShow,
	}
	return tmp.ToProduct()
}

func (r *PlatformProductRemarkBodyReq) ToProduct() ptypes.PlatformProductRemarkReq {
	return ptypes.PlatformProductRemarkReq{Remark: r.Remark}
}

func (r *ArticleUpdateBodyReq) ToContent() ctypes.ArticleSaveReq {
	return r.ArticleSaveReq.ToContent()
}

func (r *ArticleCommentPatchBodyReq) ToContent() ctypes.ArticleCommentPatchReq {
	return ctypes.ArticleCommentPatchReq{Status: r.Status}
}

func (r *ArticleCategoryUpdateBodyReq) ToContent() ctypes.ArticleCategorySaveReq {
	return ctypes.ArticleCategorySaveReq{ParentID: r.ParentID, Name: r.Name, Sort: r.Sort, Status: r.Status}
}

func (r *ArticleTopBodyReq) ToContent() ctypes.ArticleTopReq {
	return ctypes.ArticleTopReq{IsTop: r.IsTop}
}

func (r *ArticleRemarkBodyReq) ToContent() ctypes.ArticleRemarkReq {
	return ctypes.ArticleRemarkReq{Remark: r.Remark}
}

func (r *ArticleAuditBodyReq) ToContent() ctypes.ArticleAuditReq {
	return ctypes.ArticleAuditReq{Pass: r.Pass, RejectReason: r.RejectReason}
}

func (r *UserArticleUpdateBodyReq) ToContent() ctypes.ArticleSaveReq {
	return ctypes.ArticleSaveReq{
		CategoryID: r.CategoryID, Title: r.Title, CoverURL: r.CoverURL,
		Content: r.Content, ImageURLs: r.ImageURLs,
	}
}

func (r *ShopRoleUpdateBodyReq) ToShopOps() sotypes.ShopRoleReq {
	return sotypes.ShopRoleReq{Code: r.Code, Name: r.Name, Remark: r.Remark, MenuIDs: r.MenuIDs}
}

func (r *BannerUpdateBodyReq) ToBanner() (title, imageURL, linkType string, linkID uint64, sort int, status, startAt, endAt string) {
	return r.Title, r.ImageURL, r.LinkType, r.LinkID, r.Sort, r.Status, r.StartAt, r.EndAt
}

