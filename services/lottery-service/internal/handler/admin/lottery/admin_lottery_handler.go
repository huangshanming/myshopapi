package lottery

import (
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/xerr"
	"mymall/services/lottery-service/internal/biz"
	"mymall/services/lottery-service/internal/svc"
	"mymall/services/lottery-service/internal/types"
)

func ListActivitiesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		list, total, err := biz.NewLotteryLogic(svcCtx).AdminListActivities(r.Context(), req.Page, req.PageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
	}
}

func CreateActivityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ActivitySaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		vo, err := biz.NewLotteryLogic(svcCtx).AdminCreateActivity(r.Context(), biz.ActivitySaveReq{
			Title: req.Title, Status: req.Status, CostPoints: req.CostPoints,
			DailyLimit: req.DailyLimit, StartAt: req.StartAt, EndAt: req.EndAt,
		})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}

func GetActivityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		vo, err := biz.NewLotteryLogic(svcCtx).AdminGetActivity(r.Context(), req.Id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}

func UpdateActivityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ActivityUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		vo, err := biz.NewLotteryLogic(svcCtx).AdminUpdateActivity(r.Context(), req.Id, biz.ActivitySaveReq{
			Title: req.Title, Status: req.Status, CostPoints: req.CostPoints,
			DailyLimit: req.DailyLimit, StartAt: req.StartAt, EndAt: req.EndAt,
		})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}

func SavePrizesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PrizesSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		items := make([]biz.PrizeSaveItem, 0, len(req.Prizes))
		for _, p := range req.Prizes {
			items = append(items, biz.PrizeSaveItem{
				Slot: p.Slot, Name: p.Name, CoverURL: p.CoverURL, PrizeType: p.PrizeType,
				PointsAmount: p.PointsAmount, Weight: p.Weight, Stock: p.Stock, StockStrict: p.StockStrict,
			})
		}
		vo, err := biz.NewLotteryLogic(svcCtx).AdminSavePrizes(r.Context(), req.Id, items)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}

func ListRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminRecordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		list, total, err := biz.NewLotteryLogic(svcCtx).AdminListRecords(r.Context(), req.ActivityID, req.PrizeType, req.Page, req.PageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
	}
}

func ListOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminOrdersReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		list, total, err := biz.NewLotteryLogic(svcCtx).AdminListOrders(r.Context(), req.FulfillStatus, req.Page, req.PageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
	}
}

func ShipOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminShipReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := biz.NewLotteryLogic(svcCtx).AdminShipOrder(r.Context(), req.Id, req.ShipCompany, req.ShipNo); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.EmptyResp{})
	}
}

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, hdr, err := r.FormFile("file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "请选择图片文件"))
			return
		}
		defer file.Close()
		data, err := io.ReadAll(io.LimitReader(file, 6<<20))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "读取文件失败"))
			return
		}
		url, err := biz.NewLotteryLogic(svcCtx).SaveUpload(hdr.Filename, data)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.URLResp{URL: url})
	}
}
