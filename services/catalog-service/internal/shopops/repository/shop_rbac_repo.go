package repository

import (
	"context"
	"errors"
	"sort"

	"mymall/common/password"
	"mymall/services/catalog-service/internal/shopops/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	shopRoleColumns = "id, shop_id, code, name, status, remark, created_at, updated_at"
	shopMenuColumns = "id, parent_id, name, type, path, component, icon, perms, sort, visible, status, created_at, updated_at"
)

type ShopRBACRepository struct {
	conn sqlx.SqlConn
}

func NewShopRBACRepository(conn sqlx.SqlConn) *ShopRBACRepository {
	return &ShopRBACRepository{conn: conn}
}

func (r *ShopRBACRepository) EnsureOwnerRole(ctx context.Context, shopID, userID uint64) error {
	_ = r.EnsureShopMenus(ctx)
	var role model.ShopRole
	err := r.conn.QueryRowCtx(ctx, &role,
		"SELECT "+shopRoleColumns+" FROM shop_roles WHERE shop_id=? AND code=? LIMIT 1",
		shopID, "shop_owner",
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		id, insErr := lastInsertID(ctx, r.conn,
			"INSERT INTO shop_roles (shop_id, code, name, status) VALUES (?, ?, ?, ?)",
			shopID, "shop_owner", "店主", 1,
		)
		if insErr != nil {
			return insErr
		}
		role = model.ShopRole{ID: id, ShopID: shopID, Code: "shop_owner", Name: "店主", Status: 1}
	} else if err != nil {
		return err
	}
	var menus []model.ShopMenu
	_ = r.conn.QueryRowsCtx(ctx, &menus, "SELECT "+shopMenuColumns+" FROM shop_menus WHERE status=1")
	for _, m := range menus {
		cnt, _ := countCtx(ctx, r.conn,
			"SELECT COUNT(*) FROM shop_role_menus WHERE role_id=? AND menu_id=?", role.ID, m.ID)
		if cnt == 0 {
			_, _ = r.conn.ExecCtx(ctx,
				"INSERT INTO shop_role_menus (role_id, menu_id) VALUES (?, ?)", role.ID, m.ID)
		}
	}
	var ur model.ShopUserRole
	err = r.conn.QueryRowCtx(ctx, &ur,
		"SELECT shop_id, user_id, role_id FROM shop_user_roles WHERE shop_id=? AND user_id=? AND role_id=? LIMIT 1",
		shopID, userID, role.ID,
	)
	if errors.Is(err, sqlx.ErrNotFound) {
		_, err = r.conn.ExecCtx(ctx,
			"INSERT INTO shop_user_roles (shop_id, user_id, role_id) VALUES (?, ?, ?)",
			shopID, userID, role.ID,
		)
		return err
	}
	return err
}

// EnsureShopMenus 幂等写入分层商家菜单（目录 → 页面 → 按钮）
func (r *ShopRBACRepository) EnsureShopMenus(ctx context.Context) error {
	seed := []model.ShopMenu{
		{ID: 100, ParentID: 0, Name: "商品中心", Type: "dir", Path: "", Icon: "Goods", Perms: "", Sort: 10, Visible: 1, Status: 1},
		{ID: 101, ParentID: 0, Name: "库存管理", Type: "dir", Path: "", Icon: "Box", Perms: "", Sort: 20, Visible: 1, Status: 1},
		{ID: 102, ParentID: 0, Name: "订单中心", Type: "dir", Path: "", Icon: "List", Perms: "", Sort: 30, Visible: 1, Status: 1},
		{ID: 104, ParentID: 0, Name: "内容中心", Type: "dir", Path: "", Icon: "Document", Perms: "", Sort: 40, Visible: 1, Status: 1},
		{ID: 105, ParentID: 0, Name: "营销中心", Type: "dir", Path: "", Icon: "Present", Perms: "", Sort: 50, Visible: 1, Status: 1},
		{ID: 103, ParentID: 0, Name: "店铺设置", Type: "dir", Path: "", Icon: "Setting", Perms: "", Sort: 90, Visible: 1, Status: 1},

		{ID: 1, ParentID: 100, Name: "商品列表", Type: "menu", Path: "/merchant/products", Component: "merchant/Products", Icon: "Goods", Perms: "product:list", Sort: 10, Visible: 1, Status: 1},
		{ID: 2, ParentID: 100, Name: "发布商品", Type: "menu", Path: "/merchant/products/edit", Component: "merchant/ProductEdit", Icon: "Edit", Perms: "product:edit", Sort: 11, Visible: 1, Status: 1},
		{ID: 3, ParentID: 100, Name: "回收站", Type: "menu", Path: "/merchant/products/recycle", Component: "merchant/ProductRecycle", Icon: "Delete", Perms: "product:recycle", Sort: 12, Visible: 1, Status: 1},
		{ID: 7, ParentID: 100, Name: "操作日志", Type: "menu", Path: "/merchant/products/op-logs", Component: "merchant/OpLogs", Icon: "Document", Perms: "product:list", Sort: 13, Visible: 1, Status: 1},

		{ID: 4, ParentID: 101, Name: "库存预警", Type: "menu", Path: "/merchant/stocks/warnings", Component: "merchant/StockWarnings", Icon: "Warning", Perms: "stock:warn", Sort: 10, Visible: 1, Status: 1},
		{ID: 5, ParentID: 102, Name: "店铺订单", Type: "menu", Path: "/merchant/orders", Component: "merchant/Orders", Icon: "List", Perms: "order:list", Sort: 10, Visible: 1, Status: 1},
		{ID: 19, ParentID: 102, Name: "售后管理", Type: "menu", Path: "/merchant/after-sales", Component: "merchant/after-sales/AfterSales", Icon: "RefreshLeft", Perms: "aftersale:list", Sort: 11, Visible: 1, Status: 1},
		{ID: 27, ParentID: 102, Name: "评价管理", Type: "menu", Path: "/merchant/reviews", Component: "merchant/Reviews", Icon: "ChatDotRound", Perms: "product:review:list", Sort: 12, Visible: 1, Status: 1},
		{ID: 30, ParentID: 105, Name: "首页推广", Type: "menu", Path: "/merchant/homepage", Component: "merchant/HomepagePromote", Icon: "Promotion", Perms: "homepage:list", Sort: 20, Visible: 1, Status: 1},
		{ID: 32, ParentID: 105, Name: "主题坑位", Type: "menu", Path: "/merchant/themes", Component: "merchant/ThemePromote", Icon: "Grid", Perms: "theme:list", Sort: 25, Visible: 1, Status: 1},
		{ID: 34, ParentID: 105, Name: "优惠券", Type: "menu", Path: "/merchant/coupons", Component: "merchant/Coupons", Icon: "Ticket", Perms: "coupon:list", Sort: 30, Visible: 1, Status: 1},
		{ID: 8, ParentID: 104, Name: "我的文章", Type: "menu", Path: "/merchant/articles", Component: "merchant/Articles", Icon: "Document", Perms: "article:list", Sort: 10, Visible: 1, Status: 1},
		{ID: 9, ParentID: 104, Name: "发布文章", Type: "menu", Path: "/merchant/articles/edit", Component: "merchant/ArticleEdit", Icon: "Edit", Perms: "article:edit", Sort: 11, Visible: 1, Status: 1},
		{ID: 24, ParentID: 105, Name: "秒杀报名", Type: "menu", Path: "/merchant/seckill", Component: "merchant/SeckillApply", Icon: "Timer", Perms: "seckill:apply", Sort: 10, Visible: 1, Status: 1},
		{ID: 25, ParentID: 103, Name: "店铺钱包", Type: "menu", Path: "/merchant/wallet", Component: "merchant/Wallet", Icon: "Wallet", Perms: "wallet:view", Sort: 5, Visible: 1, Status: 1},
		{ID: 6, ParentID: 103, Name: "员工权限", Type: "menu", Path: "/merchant/staff", Component: "merchant/Staff", Icon: "User", Perms: "shop:staff", Sort: 10, Visible: 1, Status: 1},
		{ID: 10, ParentID: 103, Name: "消息通知", Type: "menu", Path: "/merchant/notifications", Component: "merchant/Notifications", Icon: "Bell", Perms: "notif:list", Sort: 20, Visible: 1, Status: 1},

		{ID: 11, ParentID: 1, Name: "商品新增", Type: "button", Perms: "product:add", Sort: 1, Visible: 1, Status: 1},
		{ID: 12, ParentID: 1, Name: "商品编辑", Type: "button", Perms: "product:edit", Sort: 2, Visible: 1, Status: 1},
		{ID: 13, ParentID: 1, Name: "商品上下架", Type: "button", Perms: "product:status", Sort: 3, Visible: 1, Status: 1},
		{ID: 14, ParentID: 1, Name: "批量操作", Type: "button", Perms: "product:batch", Sort: 4, Visible: 1, Status: 1},
		{ID: 15, ParentID: 1, Name: "导入导出", Type: "button", Perms: "product:import", Sort: 5, Visible: 1, Status: 1},
		{ID: 16, ParentID: 8, Name: "文章新增", Type: "button", Perms: "article:add", Sort: 1, Visible: 1, Status: 1},
		{ID: 17, ParentID: 8, Name: "文章编辑", Type: "button", Perms: "article:edit", Sort: 2, Visible: 1, Status: 1},
		{ID: 18, ParentID: 8, Name: "文章删除", Type: "button", Perms: "article:delete", Sort: 3, Visible: 1, Status: 1},
		{ID: 20, ParentID: 5, Name: "订单发货", Type: "button", Perms: "order:ship", Sort: 1, Visible: 1, Status: 1},
		{ID: 21, ParentID: 5, Name: "订单完成", Type: "button", Perms: "order:complete", Sort: 2, Visible: 1, Status: 1},
		{ID: 22, ParentID: 5, Name: "订单备注", Type: "button", Perms: "order:remark", Sort: 3, Visible: 1, Status: 1},
		{ID: 23, ParentID: 19, Name: "售后处理", Type: "button", Perms: "aftersale:handle", Sort: 1, Visible: 1, Status: 1},
		{ID: 26, ParentID: 24, Name: "提交报名", Type: "button", Perms: "seckill:apply", Sort: 1, Visible: 1, Status: 1},
		{ID: 28, ParentID: 27, Name: "评价回复", Type: "button", Perms: "product:review:reply", Sort: 1, Visible: 1, Status: 1},
		{ID: 29, ParentID: 27, Name: "评价删除", Type: "button", Perms: "product:review:delete", Sort: 2, Visible: 1, Status: 1},
		{ID: 31, ParentID: 30, Name: "购买展位", Type: "button", Perms: "homepage:buy", Sort: 1, Visible: 1, Status: 1},
		{ID: 33, ParentID: 32, Name: "购买坑位", Type: "button", Perms: "theme:buy", Sort: 1, Visible: 1, Status: 1},
		{ID: 35, ParentID: 34, Name: "编辑优惠券", Type: "button", Perms: "coupon:edit", Sort: 1, Visible: 1, Status: 1},
		{ID: 36, ParentID: 34, Name: "发放优惠券", Type: "button", Perms: "coupon:grant", Sort: 2, Visible: 1, Status: 1},
	}
	for _, m := range seed {
		var exists model.ShopMenu
		err := r.conn.QueryRowCtx(ctx, &exists, "SELECT id FROM shop_menus WHERE id=? LIMIT 1", m.ID)
		if errors.Is(err, sqlx.ErrNotFound) {
			_, err = r.conn.ExecCtx(ctx,
				`INSERT INTO shop_menus (id, parent_id, name, type, path, component, icon, perms, sort, visible, status)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				m.ID, m.ParentID, m.Name, m.Type, m.Path, m.Component, m.Icon, m.Perms, m.Sort, m.Visible, m.Status,
			)
			if err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		query, args, err := buildUpdate("shop_menus", map[string]interface{}{
			"parent_id": m.ParentID, "name": m.Name, "type": m.Type, "path": m.Path,
			"component": m.Component, "icon": m.Icon, "perms": m.Perms,
			"sort": m.Sort, "visible": m.Visible, "status": m.Status,
		}, "id=?", m.ID)
		if err != nil {
			return err
		}
		_, err = r.conn.ExecCtx(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ShopRBACRepository) ListRoleMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	type row struct {
		MenuID uint64 `db:"menu_id"`
	}
	var rows []row
	err := r.conn.QueryRowsCtx(ctx, &rows, "SELECT menu_id FROM shop_role_menus WHERE role_id=?", roleID)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, len(rows))
	for i, row := range rows {
		ids[i] = row.MenuID
	}
	return ids, nil
}

func (r *ShopRBACRepository) IsOwner(ctx context.Context, shopID, userID uint64) bool {
	cnt, _ := countCtx(ctx, r.conn,
		`SELECT COUNT(*) FROM shop_user_roles ur
		 JOIN shop_roles r ON r.id = ur.role_id
		 WHERE ur.shop_id=? AND ur.user_id=? AND r.code=?`,
		shopID, userID, "shop_owner",
	)
	return cnt > 0
}

func (r *ShopRBACRepository) HasPerm(ctx context.Context, shopID, userID uint64, code string) bool {
	if r.IsOwner(ctx, shopID, userID) {
		return true
	}
	cnt, _ := countCtx(ctx, r.conn,
		`SELECT COUNT(*) FROM shop_user_roles ur
		 JOIN shop_role_menus rm ON rm.role_id = ur.role_id
		 JOIN shop_menus m ON m.id = rm.menu_id
		 WHERE ur.shop_id=? AND ur.user_id=? AND m.perms=? AND m.status=1`,
		shopID, userID, code,
	)
	return cnt > 0
}

func (r *ShopRBACRepository) ListPerms(ctx context.Context, shopID, userID uint64) ([]string, error) {
	type row struct {
		Perms string `db:"perms"`
	}
	if r.IsOwner(ctx, shopID, userID) {
		var rows []row
		err := r.conn.QueryRowsCtx(ctx, &rows, "SELECT perms FROM shop_menus WHERE status=1 AND perms<>''")
		if err != nil {
			return nil, err
		}
		out := make([]string, len(rows))
		for i, row := range rows {
			out[i] = row.Perms
		}
		return out, nil
	}
	var rows []row
	err := r.conn.QueryRowsCtx(ctx, &rows,
		`SELECT DISTINCT m.perms FROM shop_user_roles ur
		 JOIN shop_role_menus rm ON rm.role_id = ur.role_id
		 JOIN shop_menus m ON m.id = rm.menu_id
		 WHERE ur.shop_id=? AND ur.user_id=? AND m.perms<>'' AND m.status=1`,
		shopID, userID,
	)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.Perms
	}
	return out, nil
}

func (r *ShopRBACRepository) MenuTree(ctx context.Context) ([]model.ShopMenu, error) {
	var list []model.ShopMenu
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+shopMenuColumns+" FROM shop_menus WHERE status=1 ORDER BY sort ASC, id ASC",
	)
	return list, err
}

func (r *ShopRBACRepository) ListMenusForUser(ctx context.Context, shopID, userID uint64) ([]model.ShopMenu, error) {
	if r.IsOwner(ctx, shopID, userID) {
		return r.MenuTree(ctx)
	}
	var list []model.ShopMenu
	err := r.conn.QueryRowsCtx(ctx, &list,
		`SELECT DISTINCT m.id, m.parent_id, m.name, m.type, m.path, m.component, m.icon, m.perms, m.sort, m.visible, m.status, m.created_at, m.updated_at
		 FROM shop_menus m
		 JOIN shop_role_menus rm ON rm.menu_id = m.id
		 JOIN shop_user_roles ur ON ur.role_id = rm.role_id
		 WHERE ur.shop_id=? AND ur.user_id=? AND m.status=1
		 ORDER BY m.sort ASC, m.id ASC`,
		shopID, userID,
	)
	if err != nil {
		return nil, err
	}
	return r.withAncestors(ctx, list)
}

func (r *ShopRBACRepository) withAncestors(ctx context.Context, list []model.ShopMenu) ([]model.ShopMenu, error) {
	byID := map[uint64]model.ShopMenu{}
	for _, m := range list {
		byID[m.ID] = m
	}
	need := make([]uint64, 0)
	for _, m := range list {
		pid := m.ParentID
		for pid > 0 {
			if _, ok := byID[pid]; ok {
				break
			}
			need = append(need, pid)
			var p model.ShopMenu
			err := r.conn.QueryRowCtx(ctx, &p,
				"SELECT "+shopMenuColumns+" FROM shop_menus WHERE id=? AND status=1 LIMIT 1", pid,
			)
			if errors.Is(err, sqlx.ErrNotFound) {
				break
			}
			if err != nil {
				break
			}
			byID[p.ID] = p
			pid = p.ParentID
		}
	}
	_ = need
	out := make([]model.ShopMenu, 0, len(byID))
	for _, m := range byID {
		out = append(out, m)
	}
	return out, nil
}

func BuildShopMenuTree(menus []model.ShopMenu) []map[string]interface{} {
	byParent := map[uint64][]model.ShopMenu{}
	for _, m := range menus {
		byParent[m.ParentID] = append(byParent[m.ParentID], m)
	}
	for pid := range byParent {
		sort.SliceStable(byParent[pid], func(i, j int) bool {
			if byParent[pid][i].Sort == byParent[pid][j].Sort {
				return byParent[pid][i].ID < byParent[pid][j].ID
			}
			return byParent[pid][i].Sort < byParent[pid][j].Sort
		})
	}
	var walk func(parent uint64) []map[string]interface{}
	walk = func(parent uint64) []map[string]interface{} {
		out := make([]map[string]interface{}, 0)
		for _, m := range byParent[parent] {
			node := map[string]interface{}{
				"id": m.ID, "parent_id": m.ParentID, "name": m.Name, "type": m.Type,
				"path": m.Path, "perms": m.Perms, "icon": m.Icon, "visible": m.Visible, "sort": m.Sort,
			}
			if kids := walk(m.ID); len(kids) > 0 {
				node["children"] = kids
			}
			out = append(out, node)
		}
		return out
	}
	return walk(0)
}

func (r *ShopRBACRepository) ListRoles(ctx context.Context, shopID uint64) ([]model.ShopRole, error) {
	var list []model.ShopRole
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+shopRoleColumns+" FROM shop_roles WHERE shop_id=? ORDER BY id ASC", shopID,
	)
	return list, err
}

func (r *ShopRBACRepository) SaveRole(ctx context.Context, role *model.ShopRole, menuIDs []uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if role.ID == 0 {
			id, err := lastInsertID(ctx, session,
				"INSERT INTO shop_roles (shop_id, code, name, status, remark) VALUES (?, ?, ?, ?, ?)",
				role.ShopID, role.Code, role.Name, role.Status, role.Remark,
			)
			if err != nil {
				return err
			}
			role.ID = id
		} else {
			if _, err := session.ExecCtx(ctx,
				"UPDATE shop_roles SET name=?, remark=?, status=? WHERE id=?",
				role.Name, role.Remark, role.Status, role.ID,
			); err != nil {
				return err
			}
			if _, err := execSession(ctx, session, "DELETE FROM shop_role_menus WHERE role_id=?", role.ID); err != nil {
				return err
			}
		}
		for _, mid := range menuIDs {
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO shop_role_menus (role_id, menu_id) VALUES (?, ?)", role.ID, mid,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ShopRBACRepository) BindStaff(ctx context.Context, shopID, userID, roleID uint64) error {
	_, _ = r.conn.ExecCtx(ctx, "DELETE FROM shop_user_roles WHERE shop_id=? AND user_id=?", shopID, userID)
	if _, err := r.conn.ExecCtx(ctx,
		"INSERT INTO shop_user_roles (shop_id, user_id, role_id) VALUES (?, ?, ?)",
		shopID, userID, roleID,
	); err != nil {
		return err
	}
	_, _ = r.conn.ExecCtx(ctx,
		"UPDATE users SET role='merchant_staff' WHERE id=? AND role IN ('user','')", userID,
	)
	return r.EnsureShopMember(ctx, shopID, userID, "staff")
}

// EnsureShopMember 写入店铺成员关系（登录取 shop_id 用）
func (r *ShopRBACRepository) EnsureShopMember(ctx context.Context, shopID, userID uint64, memberRole string) error {
	if shopID == 0 || userID == 0 {
		return errors.New("店铺或用户无效")
	}
	if memberRole == "" {
		memberRole = "staff"
	}
	cnt, _ := countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM shop_members WHERE shop_id=? AND user_id=?", shopID, userID,
	)
	if cnt > 0 {
		_, err := r.conn.ExecCtx(ctx,
			"UPDATE shop_members SET member_role=? WHERE shop_id=? AND user_id=?",
			memberRole, shopID, userID,
		)
		return err
	}
	_, err := r.conn.ExecCtx(ctx,
		"INSERT INTO shop_members (shop_id, user_id, member_role) VALUES (?, ?, ?)",
		shopID, userID, memberRole,
	)
	return err
}

func (r *ShopRBACRepository) FindUserIDByMobile(ctx context.Context, mobile string) (uint64, error) {
	var id uint64
	err := r.conn.QueryRowCtx(ctx, &id,
		"SELECT id FROM users WHERE mobile=? AND (deleted_at IS NULL OR deleted_at='0000-00-00 00:00:00') LIMIT 1",
		mobile,
	)
	if errors.Is(err, sqlx.ErrNotFound) || id == 0 {
		return 0, errors.New("用户不存在，请先创建账号或改用「新建店员」")
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateStaffUser 商家新建店员账号；手机号已存在时返回明确提示
func (r *ShopRBACRepository) CreateStaffUser(ctx context.Context, mobile, plainPwd, nickname string) (uint64, error) {
	if mobile == "" || len(mobile) != 11 {
		return 0, errors.New("请输入 11 位手机号")
	}
	if plainPwd == "" || len(plainPwd) < 6 {
		return 0, errors.New("密码至少 6 位")
	}
	var id uint64
	err := r.conn.QueryRowCtx(ctx, &id, "SELECT id FROM users WHERE mobile=? LIMIT 1", mobile)
	if err == nil && id > 0 {
		return id, errors.New("该手机号已注册，请改用「绑定已有账号」")
	}
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return 0, err
	}
	if nickname == "" {
		nickname = mobile
	}
	id, err = lastInsertID(ctx, r.conn,
		`INSERT INTO users (mobile, password, nickname, status, role, avatar, gender)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		mobile, password.Hash(plainPwd), nickname, 1, "merchant_staff", "", 0,
	)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("创建用户失败")
	}
	return id, nil
}

type staffRow struct {
	UserID   uint64 `db:"user_id"`
	RoleID   uint64 `db:"role_id"`
	RoleName string `db:"role_name"`
	Mobile   string `db:"mobile"`
	Nickname string `db:"nickname"`
}

func (r *ShopRBACRepository) ListStaff(ctx context.Context, shopID uint64) ([]map[string]interface{}, error) {
	var rows []staffRow
	err := r.conn.QueryRowsCtx(ctx, &rows,
		`SELECT ur.user_id, ur.role_id, r.name AS role_name, u.mobile, u.nickname
		 FROM shop_user_roles ur
		 JOIN shop_roles r ON r.id = ur.role_id
		 JOIN users u ON u.id = ur.user_id
		 WHERE ur.shop_id=?`,
		shopID,
	)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		out = append(out, map[string]interface{}{
			"user_id":   row.UserID,
			"role_id":   row.RoleID,
			"role_name": row.RoleName,
			"mobile":    row.Mobile,
			"nickname":  row.Nickname,
		})
	}
	return out, nil
}
