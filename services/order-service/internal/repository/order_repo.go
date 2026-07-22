package repository

import (
	"context"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/order-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	orderColumns = "id, order_no, user_id, shop_id, total_amount, goods_amount, discount_amount, pay_amount, user_coupon_id, receiver_name, receiver_phone, receiver_address, ship_company, ship_no, shipped_at, completed_at, reviewed_at, remark, status, created_at, updated_at"
	orderItemColumns = "id, order_id, product_id, sku_id, product_name, sku_snapshot, price, quantity, seckill_entry_id, created_at"
	afterSaleColumns = "id, order_id, order_no, user_id, shop_id, type, reason, amount, status, admin_remark, handled_by, created_at, updated_at"
)

type OrderRepository struct {
	conn sqlx.SqlConn
}

func NewOrderRepository(conn sqlx.SqlConn) *OrderRepository {
	return &OrderRepository{conn: conn}
}

type OrderListFilter struct {
	ShopID   uint64
	UserID   uint64
	Status   string
	OrderNo  string
	Page     int
	PageSize int
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order, items []model.OrderItem) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		id, err := lastInsertID(ctx, session,
			`INSERT INTO orders (order_no, user_id, shop_id, total_amount, goods_amount, discount_amount, pay_amount, user_coupon_id, receiver_name, receiver_phone, receiver_address, ship_company, ship_no, remark, status)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			order.OrderNo, order.UserID, order.ShopID, order.TotalAmount, order.GoodsAmount, order.DiscountAmount, order.PayAmount,
			order.UserCouponID, order.ReceiverName, order.ReceiverPhone, order.ReceiverAddress, order.ShipCompany, order.ShipNo, order.Remark, order.Status,
		)
		if err != nil {
			return err
		}
		order.ID = id
		for i := range items {
			items[i].OrderID = id
			itemID, err := lastInsertID(ctx, session,
				`INSERT INTO order_items (order_id, product_id, sku_id, product_name, sku_snapshot, price, quantity, seckill_entry_id)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				items[i].OrderID, items[i].ProductID, items[i].SkuID, items[i].ProductName, items[i].SkuSnapshot,
				items[i].Price, items[i].Quantity, items[i].SeckillEntryID,
			)
			if err != nil {
				return err
			}
			items[i].ID = itemID
		}
		return nil
	})
}

func (r *OrderRepository) loadItems(ctx context.Context, orders []model.Order) error {
	if len(orders) == 0 {
		return nil
	}
	ids := make([]uint64, len(orders))
	for i, o := range orders {
		ids[i] = o.ID
	}
	var items []model.OrderItem
	err := r.conn.QueryRowsCtx(ctx, &items,
		"SELECT "+orderItemColumns+" FROM order_items WHERE order_id IN ("+placeholders(len(ids))+") ORDER BY id ASC",
		inArgs(ids)...,
	)
	if err != nil {
		return err
	}
	m := map[uint64][]model.OrderItem{}
	for _, it := range items {
		m[it.OrderID] = append(m[it.OrderID], it)
	}
	for i := range orders {
		orders[i].Items = m[orders[i].ID]
	}
	return nil
}

func (r *OrderRepository) findOne(ctx context.Context, where string, args ...any) (*model.Order, error) {
	var order model.Order
	err := r.conn.QueryRowCtx(ctx, &order, "SELECT "+orderColumns+" FROM orders WHERE "+where+" LIMIT 1", args...)
	if err != nil {
		return nil, err
	}
	list := []model.Order{order}
	if err := r.loadItems(ctx, list); err != nil {
		return nil, err
	}
	return &list[0], nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id, userID uint64) (*model.Order, error) {
	return r.findOne(ctx, "id = ? AND user_id = ?", id, userID)
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return r.List(ctx, OrderListFilter{UserID: userID, Page: page, PageSize: pageSize})
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderNo, status string) error {
	_, err := r.conn.ExecCtx(ctx, "UPDATE orders SET status=? WHERE order_no=?", status, orderNo)
	return err
}

func (r *OrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	return r.findOne(ctx, "order_no = ?", orderNo)
}

func (r *OrderRepository) Cancel(ctx context.Context, orderID, userID uint64) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE orders SET status=? WHERE id=? AND user_id=? AND status IN (?, ?)",
		model.OrderStatusCancelled, orderID, userID, model.OrderStatusPending, model.OrderStatusConfirmed,
	)
	return err
}

func (r *OrderRepository) ListByShop(ctx context.Context, shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return r.List(ctx, OrderListFilter{ShopID: shopID, Page: page, PageSize: pageSize})
}

func (r *OrderRepository) ListAll(ctx context.Context, shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	return r.List(ctx, OrderListFilter{ShopID: shopID, Page: page, PageSize: pageSize})
}

func (r *OrderRepository) List(ctx context.Context, f OrderListFilter) ([]model.Order, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	where := []string{"1=1"}
	args := make([]any, 0, 8)
	if f.ShopID > 0 {
		where = append(where, "shop_id=?")
		args = append(args, f.ShopID)
	}
	if f.UserID > 0 {
		where = append(where, "user_id=?")
		args = append(args, f.UserID)
	}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if f.OrderNo != "" {
		where = append(where, "order_no LIKE ?")
		args = append(args, "%"+f.OrderNo+"%")
	}
	w := strings.Join(where, " AND ")
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM orders WHERE "+w, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), (f.Page-1)*f.PageSize, f.PageSize)
	var orders []model.Order
	err = r.conn.QueryRowsCtx(ctx, &orders,
		"SELECT "+orderColumns+" FROM orders WHERE "+w+" ORDER BY id DESC LIMIT ?, ?",
		listArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	if err := r.loadItems(ctx, orders); err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

type StatusCountRow struct {
	Status string `db:"status" json:"status"`
	Count  int64  `db:"count" json:"count"`
}

func (r *OrderRepository) CountByUserStatus(ctx context.Context, userID uint64) ([]StatusCountRow, error) {
	var rows []StatusCountRow
	err := r.conn.QueryRowsCtx(ctx, &rows,
		"SELECT status, COUNT(*) AS count FROM orders WHERE user_id=? GROUP BY status", userID,
	)
	return rows, err
}

func (r *OrderRepository) CountOpenAfterSalesByUser(ctx context.Context, userID uint64) (int64, error) {
	return countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM order_after_sales WHERE user_id=? AND status IN (?, ?)",
		userID, model.AfterSalePending, model.AfterSaleApproved,
	)
}

func (r *OrderRepository) FindByIDAndShop(ctx context.Context, id, shopID uint64) (*model.Order, error) {
	return r.findOne(ctx, "id = ? AND shop_id = ?", id, shopID)
}

func (r *OrderRepository) FindByIDAdmin(ctx context.Context, id uint64) (*model.Order, error) {
	return r.findOne(ctx, "id = ?", id)
}

func (r *OrderRepository) Ship(ctx context.Context, id uint64, shopID uint64, company, shipNo string) error {
	now := common.LocalTime(time.Now())
	q := "UPDATE orders SET status=?, ship_company=?, ship_no=?, shipped_at=? WHERE id=? AND status=?"
	args := []any{model.OrderStatusShipped, company, shipNo, &now, id, model.OrderStatusConfirmed}
	if shopID > 0 {
		q += " AND shop_id=?"
		args = append(args, shopID)
	}
	n, err := execAffected(ctx, r.conn, q, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) Complete(ctx context.Context, id uint64, shopID uint64) error {
	now := common.LocalTime(time.Now())
	q := "UPDATE orders SET status=?, completed_at=? WHERE id=? AND status=?"
	args := []any{model.OrderStatusCompleted, &now, id, model.OrderStatusShipped}
	if shopID > 0 {
		q += " AND shop_id=?"
		args = append(args, shopID)
	}
	n, err := execAffected(ctx, r.conn, q, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) ConfirmReceive(ctx context.Context, id, userID uint64) error {
	now := common.LocalTime(time.Now())
	n, err := execAffected(ctx, r.conn,
		"UPDATE orders SET status=?, completed_at=? WHERE id=? AND user_id=? AND status=?",
		model.OrderStatusCompleted, &now, id, userID, model.OrderStatusShipped,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) MarkReviewed(ctx context.Context, id, userID uint64) error {
	now := common.LocalTime(time.Now())
	n, err := execAffected(ctx, r.conn,
		"UPDATE orders SET status=?, reviewed_at=? WHERE id=? AND user_id=? AND status=?",
		model.OrderStatusReviewed, &now, id, userID, model.OrderStatusCompleted,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) UpdateRemark(ctx context.Context, id uint64, shopID uint64, remark string) error {
	q := "UPDATE orders SET remark=? WHERE id=?"
	args := []any{remark, id}
	if shopID > 0 {
		q += " AND shop_id=?"
		args = append(args, shopID)
	}
	n, err := execAffected(ctx, r.conn, q, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) EnrichOrders(ctx context.Context, orders []model.Order) {
	if len(orders) == 0 {
		return
	}
	userIDs := make([]uint64, 0, len(orders))
	shopIDs := make([]uint64, 0, len(orders))
	seenU, seenS := map[uint64]struct{}{}, map[uint64]struct{}{}
	for _, o := range orders {
		if _, ok := seenU[o.UserID]; !ok && o.UserID > 0 {
			seenU[o.UserID] = struct{}{}
			userIDs = append(userIDs, o.UserID)
		}
		if _, ok := seenS[o.ShopID]; !ok && o.ShopID > 0 {
			seenS[o.ShopID] = struct{}{}
			shopIDs = append(shopIDs, o.ShopID)
		}
	}
	userNames := r.loadUserNames(ctx, userIDs)
	shopNames := r.loadShopNames(ctx, shopIDs)
	for i := range orders {
		orders[i].UserName = userNames[orders[i].UserID]
		orders[i].ShopName = shopNames[orders[i].ShopID]
	}
}

func (r *OrderRepository) EnrichOrder(ctx context.Context, order *model.Order) {
	if order == nil {
		return
	}
	names := r.loadUserNames(ctx, []uint64{order.UserID})
	shops := r.loadShopNames(ctx, []uint64{order.ShopID})
	order.UserName = names[order.UserID]
	order.ShopName = shops[order.ShopID]
}

func (r *OrderRepository) loadUserNames(ctx context.Context, ids []uint64) map[uint64]string {
	out := map[uint64]string{}
	if len(ids) == 0 {
		return out
	}
	type row struct {
		ID       uint64 `db:"id"`
		Nickname string `db:"nickname"`
		Mobile   string `db:"mobile"`
	}
	var rows []row
	if err := r.conn.QueryRowsCtx(ctx, &rows,
		"SELECT id, nickname, mobile FROM users WHERE id IN ("+placeholders(len(ids))+")",
		inArgs(ids)...,
	); err != nil {
		return out
	}
	for _, u := range rows {
		if u.Nickname != "" {
			out[u.ID] = u.Nickname
		} else if u.Mobile != "" {
			out[u.ID] = u.Mobile
		}
	}
	return out
}

func (r *OrderRepository) loadShopNames(ctx context.Context, ids []uint64) map[uint64]string {
	out := map[uint64]string{}
	if len(ids) == 0 {
		return out
	}
	type row struct {
		ID   uint64 `db:"id"`
		Name string `db:"name"`
	}
	var rows []row
	if err := r.conn.QueryRowsCtx(ctx, &rows,
		"SELECT id, name FROM shops WHERE id IN ("+placeholders(len(ids))+")",
		inArgs(ids)...,
	); err != nil {
		return out
	}
	for _, s := range rows {
		out[s.ID] = s.Name
	}
	return out
}

type AfterSaleListFilter struct {
	ShopID   uint64
	UserID   uint64
	Status   string
	OrderNo  string
	Page     int
	PageSize int
}

func (r *OrderRepository) CreateAfterSale(ctx context.Context, as *model.OrderAfterSale) error {
	id, err := lastInsertID(ctx, r.conn,
		`INSERT INTO order_after_sales (order_id, order_no, user_id, shop_id, type, reason, amount, status, admin_remark, handled_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		as.OrderID, as.OrderNo, as.UserID, as.ShopID, as.Type, as.Reason, as.Amount, as.Status, as.AdminRemark, as.HandledBy,
	)
	if err != nil {
		return err
	}
	as.ID = id
	return nil
}

func (r *OrderRepository) FindAfterSale(ctx context.Context, id uint64) (*model.OrderAfterSale, error) {
	var as model.OrderAfterSale
	err := r.conn.QueryRowCtx(ctx, &as, "SELECT "+afterSaleColumns+" FROM order_after_sales WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &as, nil
}

func (r *OrderRepository) ListAfterSales(ctx context.Context, f AfterSaleListFilter) ([]model.OrderAfterSale, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	where := []string{"1=1"}
	args := make([]any, 0, 8)
	if f.ShopID > 0 {
		where = append(where, "shop_id=?")
		args = append(args, f.ShopID)
	}
	if f.UserID > 0 {
		where = append(where, "user_id=?")
		args = append(args, f.UserID)
	}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if f.OrderNo != "" {
		where = append(where, "order_no LIKE ?")
		args = append(args, "%"+f.OrderNo+"%")
	}
	w := strings.Join(where, " AND ")
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM order_after_sales WHERE "+w, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), (f.Page-1)*f.PageSize, f.PageSize)
	var list []model.OrderAfterSale
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+afterSaleColumns+" FROM order_after_sales WHERE "+w+" ORDER BY id DESC LIMIT ?, ?",
		listArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	r.enrichAfterSales(ctx, list)
	return list, total, nil
}

func (r *OrderRepository) ListAfterSalesByOrder(ctx context.Context, orderID uint64) ([]model.OrderAfterSale, error) {
	var list []model.OrderAfterSale
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+afterSaleColumns+" FROM order_after_sales WHERE order_id=? ORDER BY id DESC", orderID,
	)
	if err != nil {
		return nil, err
	}
	r.enrichAfterSales(ctx, list)
	return list, nil
}

func (r *OrderRepository) HandleAfterSale(ctx context.Context, id, shopID, handledBy uint64, status, adminRemark string) error {
	q := "UPDATE order_after_sales SET status=?, admin_remark=?, handled_by=? WHERE id=?"
	args := []any{status, adminRemark, handledBy, id}
	if shopID > 0 {
		q += " AND shop_id=?"
		args = append(args, shopID)
	}
	n, err := execAffected(ctx, r.conn, q, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) CountOpenAfterSales(ctx context.Context, orderID uint64) (int64, error) {
	return countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM order_after_sales WHERE order_id=? AND status IN (?, ?)",
		orderID, model.AfterSalePending, model.AfterSaleApproved,
	)
}

func (r *OrderRepository) enrichAfterSales(ctx context.Context, list []model.OrderAfterSale) {
	if len(list) == 0 {
		return
	}
	userIDs, shopIDs := []uint64{}, []uint64{}
	seenU, seenS := map[uint64]struct{}{}, map[uint64]struct{}{}
	for _, a := range list {
		if _, ok := seenU[a.UserID]; !ok {
			seenU[a.UserID] = struct{}{}
			userIDs = append(userIDs, a.UserID)
		}
		if _, ok := seenS[a.ShopID]; !ok && a.ShopID > 0 {
			seenS[a.ShopID] = struct{}{}
			shopIDs = append(shopIDs, a.ShopID)
		}
	}
	un := r.loadUserNames(ctx, userIDs)
	sn := r.loadShopNames(ctx, shopIDs)
	for i := range list {
		list[i].UserName = un[list[i].UserID]
		list[i].ShopName = sn[list[i].ShopID]
	}
}
