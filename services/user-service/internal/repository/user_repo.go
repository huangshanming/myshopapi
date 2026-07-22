package repository

import (
	"context"
	"errors"

	"mymall/common/password"
	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const userColumns = "id, created_at, updated_at, mobile, IFNULL(password,'') AS password, IFNULL(nickname,'') AS nickname, IFNULL(avatar,'') AS avatar, gender, status, IFNULL(role,'') AS role, last_login_time, deleted_at"

type UserRepository struct {
	conn sqlx.SqlConn
}

func NewUserRepository(conn sqlx.SqlConn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) HashPassword(ctx context.Context, plain string) string {
	return password.Hash(plain)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uint64, plain string) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE users SET password=? WHERE id=? AND deleted_at IS NULL",
		password.Hash(plain), id,
	)
	return err
}

func (r *UserRepository) CreateAdmin(ctx context.Context, mobile, plain, nickname string) (*model.User, error) {
	var existing model.User
	err := r.conn.QueryRowPartialCtx(ctx, &existing,
		"SELECT "+userColumns+" FROM users WHERE mobile=? AND deleted_at IS NULL LIMIT 1", mobile,
	)
	if err == nil {
		return nil, errors.New("用户已存在")
	}
	if !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}
	if nickname == "" {
		nickname = mobile
	}
	res, err := r.conn.ExecCtx(ctx,
		"INSERT INTO users (mobile, password, nickname, status, role) VALUES (?,?,?,?,?)",
		mobile, password.Hash(plain), nickname, 1, "platform_admin",
	)
	if err != nil {
		return nil, err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *UserRepository) FindByMobile(ctx context.Context, mobile string) (*model.User, error) {
	var user model.User
	err := r.conn.QueryRowPartialCtx(ctx, &user,
		"SELECT "+userColumns+" FROM users WHERE mobile=? AND deleted_at IS NULL LIMIT 1", mobile,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) VerifyLogin(ctx context.Context, mobile, plain string) (*model.User, error) {
	user, err := r.FindByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	if user.Password != password.Hash(plain) {
		return nil, sqlx.ErrNotFound
	}
	if user.Status != 1 {
		return nil, errors.New("账号已禁用")
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, mobile, plain string) (*model.User, error) {
	var existing model.User
	err := r.conn.QueryRowPartialCtx(ctx, &existing,
		"SELECT "+userColumns+" FROM users WHERE mobile=? AND deleted_at IS NULL LIMIT 1", mobile,
	)
	if err == nil {
		return nil, errors.New("用户已存在")
	}
	if !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}

	res, err := r.conn.ExecCtx(ctx,
		"INSERT INTO users (mobile, password, nickname, status, role) VALUES (?,?,?,?,?)",
		mobile, password.Hash(plain), mobile, 1, "user",
	)
	if err != nil {
		return nil, err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id uint64, nickname, avatar, mobile string, gender int) error {
	_, err := r.conn.ExecCtx(ctx,
		"UPDATE users SET nickname=?, avatar=?, gender=?, mobile=? WHERE id=? AND deleted_at IS NULL",
		nickname, avatar, gender, mobile, id,
	)
	return err
}

func (r *UserRepository) MobileTakenByOther(ctx context.Context, mobile string, excludeID uint64) bool {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM users WHERE mobile=? AND id<>? AND deleted_at IS NULL",
		mobile, excludeID,
	)
	return err == nil && n > 0
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.conn.QueryRowPartialCtx(ctx, &user,
		"SELECT "+userColumns+" FROM users WHERE id=? AND deleted_at IS NULL LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FirstShopID 取用户所属第一家店铺（优先 shop_members，兼容 shop_user_roles）
func (r *UserRepository) FirstShopID(ctx context.Context, userID uint64) uint64 {
	var shopID uint64
	_ = r.conn.QueryRowPartialCtx(ctx, &shopID,
		"SELECT shop_id FROM shop_members WHERE user_id=? ORDER BY id ASC LIMIT 1", userID,
	)
	if shopID > 0 {
		return shopID
	}
	_ = r.conn.QueryRowPartialCtx(ctx, &shopID,
		"SELECT shop_id FROM shop_user_roles WHERE user_id=? ORDER BY shop_id ASC LIMIT 1", userID,
	)
	return shopID
}
