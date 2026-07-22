package repository

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const addressColumns = "id, user_id, receiver_name, receiver_phone, province, city, district, detail, province_code, city_code, district_code, is_default, created_at, updated_at"

func (r *UserRepository) ListAddresses(ctx context.Context, userID uint64) ([]model.UserAddress, error) {
	var list []model.UserAddress
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+addressColumns+" FROM user_addresses WHERE user_id=? ORDER BY is_default DESC, id DESC",
		userID,
	)
	return list, err
}

func (r *UserRepository) GetAddress(ctx context.Context, userID, id uint64) (*model.UserAddress, error) {
	var a model.UserAddress
	err := r.conn.QueryRowCtx(ctx, &a,
		"SELECT "+addressColumns+" FROM user_addresses WHERE id=? AND user_id=? LIMIT 1", id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *UserRepository) CreateAddress(ctx context.Context, a *model.UserAddress) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if a.IsDefault == 1 {
			if _, err := session.ExecCtx(ctx,
				"UPDATE user_addresses SET is_default=0 WHERE user_id=?", a.UserID,
			); err != nil {
				return err
			}
		} else {
			n, err := countQuery(ctx, session,
				"SELECT COUNT(*) FROM user_addresses WHERE user_id=?", a.UserID,
			)
			if err != nil {
				return err
			}
			if n == 0 {
				a.IsDefault = 1
			}
		}
		res, err := session.ExecCtx(ctx,
			`INSERT INTO user_addresses (user_id, receiver_name, receiver_phone, province, city, district, detail, province_code, city_code, district_code, is_default)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			a.UserID, a.ReceiverName, a.ReceiverPhone, a.Province, a.City, a.District, a.Detail,
			a.ProvinceCode, a.CityCode, a.DistrictCode, a.IsDefault,
		)
		if err != nil {
			return err
		}
		id, err := lastInsertID(res)
		if err != nil {
			return err
		}
		a.ID = id
		return nil
	})
}

func (r *UserRepository) UpdateAddress(ctx context.Context, userID, id uint64, a *model.UserAddress) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var existing model.UserAddress
		if err := session.QueryRowCtx(ctx, &existing,
			"SELECT "+addressColumns+" FROM user_addresses WHERE id=? AND user_id=? LIMIT 1", id, userID,
		); err != nil {
			return err
		}
		if a.IsDefault == 1 {
			if _, err := session.ExecCtx(ctx,
				"UPDATE user_addresses SET is_default=0 WHERE user_id=?", userID,
			); err != nil {
				return err
			}
		}
		_, err := session.ExecCtx(ctx,
			`UPDATE user_addresses SET receiver_name=?, receiver_phone=?, province=?, city=?, district=?, detail=?,
			 province_code=?, city_code=?, district_code=?, is_default=? WHERE id=? AND user_id=?`,
			a.ReceiverName, a.ReceiverPhone, a.Province, a.City, a.District, a.Detail,
			a.ProvinceCode, a.CityCode, a.DistrictCode, a.IsDefault, id, userID,
		)
		return err
	})
}

func (r *UserRepository) DeleteAddress(ctx context.Context, userID, id uint64) error {
	n, err := execRows(ctx, r.conn,
		"DELETE FROM user_addresses WHERE id=? AND user_id=?", id, userID,
	)
	if err != nil {
		return err
	}
	if n == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

func (r *UserRepository) SetDefaultAddress(ctx context.Context, userID, id uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var existing model.UserAddress
		if err := session.QueryRowCtx(ctx, &existing,
			"SELECT "+addressColumns+" FROM user_addresses WHERE id=? AND user_id=? LIMIT 1", id, userID,
		); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE user_addresses SET is_default=0 WHERE user_id=?", userID,
		); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx,
			"UPDATE user_addresses SET is_default=1 WHERE id=? AND user_id=?", id, userID,
		)
		return err
	})
}

func (r *UserRepository) GetAddressByID(ctx context.Context, userID, id uint64) (*model.UserAddress, error) {
	a, err := r.GetAddress(ctx, userID, id)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errors.New("收货地址不存在")
		}
		return nil, err
	}
	return a, nil
}
