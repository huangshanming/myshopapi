package biz

import (
	"context"
	"errors"
	"strings"

	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/repository"
	"mymall/services/user-service/internal/svc"
)

type PointsOrderLogic struct {
	svcCtx *svc.ServiceContext
}

func NewPointsOrderLogic(svcCtx *svc.ServiceContext) *PointsOrderLogic {
	return &PointsOrderLogic{svcCtx: svcCtx}
}

type PointsOrderVO struct {
	model.PointsExchangeOrder
	UserName   string `json:"user_name"`
	UserMobile string `json:"user_mobile"`
}

type ExchangeReq struct {
	ProductID       uint64 `json:"product_id"`
	Quantity        int    `json:"quantity"`
	ReceiverName    string `json:"receiver_name"`
	ReceiverPhone   string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
}

type ShipReq struct {
	ShipCompany string `json:"ship_company"`
	ShipNo      string `json:"ship_no"`
}

type RemarkReq struct {
	AdminRemark string `json:"admin_remark"`
}

func (l *PointsOrderLogic) enrich(ctx context.Context, o *model.PointsExchangeOrder) PointsOrderVO {
	vo := PointsOrderVO{PointsExchangeOrder: *o}
	briefs := l.svcCtx.PointsOrders.MapUserBriefs(ctx, []uint64{o.UserID})
	if b, ok := briefs[o.UserID]; ok {
		vo.UserName = b[0]
		vo.UserMobile = b[1]
	}
	return vo
}

func (l *PointsOrderLogic) AdminList(ctx context.Context, page, pageSize int, status, orderNo, keyword string, userID uint64) ([]PointsOrderVO, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	list, total, err := l.svcCtx.PointsOrders.List(ctx, page, pageSize, repository.PointsOrderListFilter{
		Status: status, OrderNo: orderNo, Keyword: keyword, UserID: userID,
	})
	if err != nil {
		return nil, 0, err
	}
	ids := make([]uint64, 0, len(list))
	seen := map[uint64]struct{}{}
	for _, o := range list {
		if _, ok := seen[o.UserID]; !ok {
			seen[o.UserID] = struct{}{}
			ids = append(ids, o.UserID)
		}
	}
	briefs := l.svcCtx.PointsOrders.MapUserBriefs(ctx, ids)
	out := make([]PointsOrderVO, 0, len(list))
	for i := range list {
		vo := PointsOrderVO{PointsExchangeOrder: list[i]}
		if b, ok := briefs[list[i].UserID]; ok {
			vo.UserName, vo.UserMobile = b[0], b[1]
		}
		out = append(out, vo)
	}
	return out, total, nil
}

func (l *PointsOrderLogic) AdminGet(ctx context.Context, id uint64) (*PointsOrderVO, error) {
	o, err := l.svcCtx.PointsOrders.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	vo := l.enrich(ctx, o)
	return &vo, nil
}

func (l *PointsOrderLogic) AdminShip(ctx context.Context, id uint64, req ShipReq) (*PointsOrderVO, error) {
	if strings.TrimSpace(req.ShipNo) == "" {
		return nil, errors.New("请填写物流单号")
	}
	o, err := l.svcCtx.PointsOrders.Ship(ctx, id, req.ShipCompany, req.ShipNo)
	if err != nil {
		return nil, err
	}
	vo := l.enrich(ctx, o)
	return &vo, nil
}

func (l *PointsOrderLogic) AdminComplete(ctx context.Context, id uint64) (*PointsOrderVO, error) {
	o, err := l.svcCtx.PointsOrders.Complete(ctx, id)
	if err != nil {
		return nil, err
	}
	vo := l.enrich(ctx, o)
	return &vo, nil
}

func (l *PointsOrderLogic) AdminCancel(ctx context.Context, id uint64, remark string) (*PointsOrderVO, error) {
	o, err := l.svcCtx.PointsOrders.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if o.Status != model.PointsOrderPending {
		return nil, errors.New("仅待发货订单可取消退积分")
	}
	tasks := NewTaskLogic(l.svcCtx)
	if _, err := tasks.RefundPoints(ctx, PointsLedgerReq{
		UserID: o.UserID, Points: o.PointsCost, ChangeType: model.PointChangeMallRefund,
		Remark: "兑换取消退回：" + o.ProductName, RefType: model.PointsOrderRefType, RefID: o.ID,
	}); err != nil {
		return nil, err
	}
	o, err = l.svcCtx.PointsOrders.CancelLocal(ctx, id, remark)
	if err != nil {
		return nil, err
	}
	vo := l.enrich(ctx, o)
	return &vo, nil
}

func (l *PointsOrderLogic) AdminRemark(ctx context.Context, id uint64, remark string) (*PointsOrderVO, error) {
	o, err := l.svcCtx.PointsOrders.Remark(ctx, id, remark)
	if err != nil {
		return nil, err
	}
	vo := l.enrich(ctx, o)
	return &vo, nil
}

func (l *PointsOrderLogic) UserExchange(ctx context.Context, userID uint64, req ExchangeReq) (*model.PointsExchangeOrder, error) {
	if req.ProductID == 0 {
		return nil, errors.New("请选择商品")
	}
	qty := req.Quantity
	if qty < 1 {
		qty = 1
	}
	o, err := l.svcCtx.PointsOrders.CreateExchangeLocal(ctx, userID, req.ProductID, qty, req.ReceiverName, req.ReceiverPhone, req.ReceiverAddress)
	if err != nil {
		return nil, err
	}
	tasks := NewTaskLogic(l.svcCtx)
	if _, err := tasks.DeductPoints(ctx, PointsLedgerReq{
		UserID: userID, Points: o.PointsCost, ChangeType: model.PointChangeMallExchange,
		Remark: "积分兑换：" + o.ProductName, RefType: model.PointsOrderRefType, RefID: o.ID,
	}); err != nil {
		_ = l.svcCtx.PointsOrders.AbortExchange(ctx, o.ID)
		return nil, err
	}
	return o, nil
}

func (l *PointsOrderLogic) UserList(ctx context.Context, userID uint64, page, pageSize int) ([]model.PointsExchangeOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return l.svcCtx.PointsOrders.List(ctx, page, pageSize, repository.PointsOrderListFilter{UserID: userID})
}

func (l *PointsOrderLogic) UserGet(ctx context.Context, userID, id uint64) (*model.PointsExchangeOrder, error) {
	o, err := l.svcCtx.PointsOrders.GetByID(ctx, id)
	if err != nil || o.UserID != userID {
		return nil, errors.New("订单不存在")
	}
	return o, nil
}
