package product

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/merchant/product"
	"mymall/services/catalog-service/internal/svc"
)

func AdjustStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewAdjustStockLogic(r.Context(), svcCtx)
		l.AdjustStock(w, r)
	}
}

func BatchStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewBatchStockLogic(r.Context(), svcCtx)
		l.BatchStock(w, r)
	}
}

func CancelScheduleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewCancelScheduleLogic(r.Context(), svcCtx)
		l.CancelSchedule(w, r)
	}
}

func DeleteAttrTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewDeleteAttrTemplateLogic(r.Context(), svcCtx)
		l.DeleteAttrTemplate(w, r)
	}
}

func DeleteTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewDeleteTagLogic(r.Context(), svcCtx)
		l.DeleteTag(w, r)
	}
}

func JobStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewJobStatusLogic(r.Context(), svcCtx)
		l.JobStatus(w, r)
	}
}

func ListAttrTemplatesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewListAttrTemplatesLogic(r.Context(), svcCtx)
		l.ListAttrTemplates(w, r)
	}
}

func ListTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewListTagsLogic(r.Context(), svcCtx)
		l.ListTags(w, r)
	}
}

func MerchantBatchProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantBatchProductsLogic(r.Context(), svcCtx)
		l.MerchantBatchProducts(w, r)
	}
}

func MerchantCopyProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantCopyProductLogic(r.Context(), svcCtx)
		l.MerchantCopyProduct(w, r)
	}
}

func MerchantCreateAttrTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantCreateAttrTemplateLogic(r.Context(), svcCtx)
		l.MerchantCreateAttrTemplate(w, r)
	}
}

func MerchantCreateProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantCreateProductLogic(r.Context(), svcCtx)
		l.MerchantCreateProduct(w, r)
	}
}

func MerchantCreateTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantCreateTagLogic(r.Context(), svcCtx)
		l.MerchantCreateTag(w, r)
	}
}

func MerchantExportProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantExportProductsLogic(r.Context(), svcCtx)
		l.MerchantExportProducts(w, r)
	}
}

func MerchantGetProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantGetProductLogic(r.Context(), svcCtx)
		l.MerchantGetProduct(w, r)
	}
}

func MerchantImportProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantImportProductsLogic(r.Context(), svcCtx)
		l.MerchantImportProducts(w, r)
	}
}

func MerchantListProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantListProductsLogic(r.Context(), svcCtx)
		l.MerchantListProducts(w, r)
	}
}

func MerchantPurgeProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantPurgeProductsLogic(r.Context(), svcCtx)
		l.MerchantPurgeProducts(w, r)
	}
}

func MerchantRestoreProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantRestoreProductsLogic(r.Context(), svcCtx)
		l.MerchantRestoreProducts(w, r)
	}
}

func MerchantScheduleProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantScheduleProductLogic(r.Context(), svcCtx)
		l.MerchantScheduleProduct(w, r)
	}
}

func MerchantSetProductStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantSetProductStatusLogic(r.Context(), svcCtx)
		l.MerchantSetProductStatus(w, r)
	}
}

func MerchantUpdateAttrTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantUpdateAttrTemplateLogic(r.Context(), svcCtx)
		l.MerchantUpdateAttrTemplate(w, r)
	}
}

func MerchantUpdateProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantUpdateProductLogic(r.Context(), svcCtx)
		l.MerchantUpdateProduct(w, r)
	}
}

func MerchantUpdateTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantUpdateTagLogic(r.Context(), svcCtx)
		l.MerchantUpdateTag(w, r)
	}
}

func MerchantUploadImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewMerchantUploadImageLogic(r.Context(), svcCtx)
		l.MerchantUploadImage(w, r)
	}
}

func OpLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewOpLogsLogic(r.Context(), svcCtx)
		l.OpLogs(w, r)
	}
}

func StockWarningsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewStockWarningsLogic(r.Context(), svcCtx)
		l.StockWarnings(w, r)
	}
}
