package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/service"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	svc *service.MerchantService
}

func NewMerchantHandler(svc *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{svc: svc}
}

func (h *MerchantHandler) Apply(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	var in service.ApplyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	app, err := h.svc.Apply(userID, in)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, app, "提交成功")
}

func (h *MerchantHandler) MyShops(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	shops, err := h.svc.MyShops(userID)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, shops, "查询成功")
}

func (h *MerchantHandler) UpdateMyShop(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	shopID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "店铺ID无效", http.StatusBadRequest)
		return
	}
	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateMyShop(shopID, userID, &shop); err != nil {
		response.Error(c, err.Error(), http.StatusForbidden)
		return
	}
	response.Success(c, nil, "更新成功")
}

func (h *MerchantHandler) AdminListApplications(c *gin.Context) {
	var page pagination.PageReq
	_ = c.ShouldBindQuery(&page)
	p, ps, _ := pagination.Normalize(&page)
	list, total, err := h.svc.ListApplications(c.Query("status"), p, ps)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{"total": total, "list": list}, "查询成功")
}

func (h *MerchantHandler) AdminApprove(c *gin.Context) {
	adminID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "申请ID无效", http.StatusBadRequest)
		return
	}
	shop, err := h.svc.Approve(appID, adminID)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, shop, "审核通过")
}

type rejectReq struct {
	Reason string `json:"reason"`
}

func (h *MerchantHandler) AdminReject(c *gin.Context) {
	adminID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "申请ID无效", http.StatusBadRequest)
		return
	}
	var req rejectReq
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.Reject(appID, adminID, req.Reason); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "已拒绝")
}

func (h *MerchantHandler) AdminListShops(c *gin.Context) {
	var page pagination.PageReq
	_ = c.ShouldBindQuery(&page)
	p, ps, _ := pagination.Normalize(&page)
	list, total, err := h.svc.ListShops(c.Query("status"), p, ps)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{"total": total, "list": list}, "查询成功")
}

func (h *MerchantHandler) AdminGetShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "店铺ID无效", http.StatusBadRequest)
		return
	}
	shop, err := h.svc.GetShop(id)
	if err != nil {
		response.Error(c, "店铺不存在", http.StatusNotFound)
		return
	}
	response.Success(c, shop, "查询成功")
}

func (h *MerchantHandler) AdminDisableShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "店铺ID无效", http.StatusBadRequest)
		return
	}
	var req rejectReq
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.DisableShop(id, req.Reason); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "已禁用")
}

func (h *MerchantHandler) AdminEnableShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "店铺ID无效", http.StatusBadRequest)
		return
	}
	if err := h.svc.EnableShop(id); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "已启用")
}
