package repository

import (
	"context"
	"errors"

	"mymall/common/password"
	"mymall/services/merchant-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	shopColumns = "id, name, logo, contact_name, contact_phone, description, category, province, city, district, address, business_license_no, legal_person, license_image, storefront_image, owner_user_id, status, reject_reason, created_at, updated_at"
	shopAppColumns = "id, user_id, shop_name, contact_name, contact_phone, description, category, province, city, district, address, business_license_no, legal_person, license_image, storefront_image, status, reject_reason, reviewed_by, reviewed_at, shop_id, created_at, updated_at"
)

type MerchantRepository struct {
	conn sqlx.SqlConn
}

func NewMerchantRepository(conn sqlx.SqlConn) *MerchantRepository {
	return &MerchantRepository{conn: conn}
}

func (r *MerchantRepository) CreateApplication(ctx context.Context, app *model.ShopApplication) error {
	res, err := r.conn.ExecCtx(ctx,
		`INSERT INTO shop_applications (user_id, shop_name, contact_name, contact_phone, description, category, province, city, district, address, business_license_no, legal_person, license_image, storefront_image, status)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		app.UserID, app.ShopName, app.ContactName, app.ContactPhone, app.Description, app.Category,
		app.Province, app.City, app.District, app.Address, app.BusinessLicenseNo, app.LegalPerson,
		app.LicenseImage, app.StorefrontImage, app.Status,
	)
	if err != nil {
		return err
	}
	id, err := lastInsertID(res)
	if err != nil {
		return err
	}
	app.ID = id
	return nil
}

func (r *MerchantRepository) FindPendingAppByUser(ctx context.Context, userID uint64) (*model.ShopApplication, error) {
	var app model.ShopApplication
	err := r.conn.QueryRowCtx(ctx, &app,
		"SELECT "+shopAppColumns+" FROM shop_applications WHERE user_id=? AND status=? LIMIT 1",
		userID, model.AppPending,
	)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *MerchantRepository) ListApplications(ctx context.Context, status string, page, pageSize int) ([]model.ShopApplication, int64, error) {
	where := "1=1"
	args := make([]any, 0, 1)
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM shop_applications WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.ShopApplication
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+shopAppColumns+" FROM shop_applications WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) FindApplication(ctx context.Context, id uint64) (*model.ShopApplication, error) {
	var app model.ShopApplication
	err := r.conn.QueryRowCtx(ctx, &app,
		"SELECT "+shopAppColumns+" FROM shop_applications WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *MerchantRepository) ApproveApplication(ctx context.Context, appID, adminID uint64) (*model.Shop, error) {
	var shop *model.Shop
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var app model.ShopApplication
		if err := session.QueryRowCtx(ctx, &app,
			"SELECT "+shopAppColumns+" FROM shop_applications WHERE id=? LIMIT 1", appID,
		); err != nil {
			return err
		}
		if app.Status != model.AppPending {
			return ErrInvalidData
		}
		shop = &model.Shop{
			Name:              app.ShopName,
			ContactName:       app.ContactName,
			ContactPhone:      app.ContactPhone,
			Description:       app.Description,
			Category:          app.Category,
			Province:          app.Province,
			City:              app.City,
			District:          app.District,
			Address:           app.Address,
			BusinessLicenseNo: app.BusinessLicenseNo,
			LegalPerson:       app.LegalPerson,
			LicenseImage:      app.LicenseImage,
			StorefrontImage:   app.StorefrontImage,
			OwnerUserID:       app.UserID,
			Status:            model.ShopApproved,
		}
		res, err := session.ExecCtx(ctx,
			`INSERT INTO shops (name, contact_name, contact_phone, description, category, province, city, district, address, business_license_no, legal_person, license_image, storefront_image, owner_user_id, status)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			shop.Name, shop.ContactName, shop.ContactPhone, shop.Description, shop.Category,
			shop.Province, shop.City, shop.District, shop.Address, shop.BusinessLicenseNo,
			shop.LegalPerson, shop.LicenseImage, shop.StorefrontImage, shop.OwnerUserID, shop.Status,
		)
		if err != nil {
			return err
		}
		shopID, err := lastInsertID(res)
		if err != nil {
			return err
		}
		shop.ID = shopID
		if _, err := session.ExecCtx(ctx,
			"INSERT INTO shop_members (shop_id, user_id, member_role) VALUES (?,?,?)",
			shop.ID, app.UserID, model.MemberOwner,
		); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx,
			"UPDATE users SET role=? WHERE id=?", "merchant_owner", app.UserID,
		); err != nil {
			return err
		}
		_, err = session.ExecCtx(ctx,
			`UPDATE shop_applications SET status=?, reviewed_by=?, shop_id=?, reviewed_at=NOW() WHERE id=?`,
			model.AppApproved, adminID, shop.ID, appID,
		)
		return err
	})
	return shop, err
}

func (r *MerchantRepository) RejectApplication(ctx context.Context, appID, adminID uint64, reason string) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE shop_applications SET status=?, reject_reason=?, reviewed_by=?, reviewed_at=NOW() WHERE id=? AND status=?`,
		model.AppRejected, reason, adminID, appID, model.AppPending,
	)
	return err
}

func (r *MerchantRepository) ListShops(ctx context.Context, status, name string, page, pageSize int) ([]model.Shop, int64, error) {
	where := "1=1"
	args := make([]any, 0, 2)
	if status != "" {
		where += " AND status=?"
		args = append(args, status)
	}
	if name != "" {
		where += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM shops WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.Shop
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+shopColumns+" FROM shops WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) ListPublicShops(ctx context.Context, page, pageSize int) ([]model.Shop, int64, error) {
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM shops WHERE status=?", model.ShopApproved)
	if err != nil {
		return nil, 0, err
	}
	var list []model.Shop
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+shopColumns+" FROM shops WHERE status=? ORDER BY id ASC LIMIT ? OFFSET ?",
		model.ShopApproved, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *MerchantRepository) FindShop(ctx context.Context, id uint64) (*model.Shop, error) {
	var shop model.Shop
	err := r.conn.QueryRowCtx(ctx, &shop,
		"SELECT "+shopColumns+" FROM shops WHERE id=? LIMIT 1", id,
	)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *MerchantRepository) UpdateShopStatus(ctx context.Context, id uint64, status, reason string) error {
	if reason != "" {
		_, err := r.conn.ExecCtx(ctx,
			"UPDATE shops SET status=?, reject_reason=? WHERE id=?", status, reason, id,
		)
		return err
	}
	_, err := r.conn.ExecCtx(ctx, "UPDATE shops SET status=? WHERE id=?", status, id)
	return err
}

func (r *MerchantRepository) UpdateShop(ctx context.Context, shop *model.Shop) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE shops SET name=?, logo=?, contact_name=?, contact_phone=?, description=?, category=?, province=?, city=?, district=?, address=?, business_license_no=?, legal_person=?, license_image=?, storefront_image=? WHERE id=?`,
		shop.Name, shop.Logo, shop.ContactName, shop.ContactPhone, shop.Description, shop.Category,
		shop.Province, shop.City, shop.District, shop.Address, shop.BusinessLicenseNo, shop.LegalPerson,
		shop.LicenseImage, shop.StorefrontImage, shop.ID,
	)
	return err
}

func (r *MerchantRepository) UpdateShopDisplay(ctx context.Context, shop *model.Shop) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE shops SET name=?, logo=?, contact_name=?, contact_phone=?, description=?, category=?, province=?, city=?, district=?, address=?, storefront_image=? WHERE id=?`,
		shop.Name, shop.Logo, shop.ContactName, shop.ContactPhone, shop.Description, shop.Category,
		shop.Province, shop.City, shop.District, shop.Address, shop.StorefrontImage, shop.ID,
	)
	return err
}

func (r *MerchantRepository) CreateShopWithOwner(ctx context.Context, shop *model.Shop, mobile, plainPwd, nickname string) (*model.Shop, error) {
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var ownerID uint64
		err := session.QueryRowCtx(ctx, &ownerID,
			"SELECT id FROM users WHERE mobile=? LIMIT 1", mobile,
		)
		if err == nil {
			if _, err := session.ExecCtx(ctx,
				"UPDATE users SET role=? WHERE id=?", "merchant_owner", ownerID,
			); err != nil {
				return err
			}
		} else if errors.Is(err, sqlx.ErrNotFound) {
			if plainPwd == "" {
				return errors.New("新用户必须设置密码")
			}
			if nickname == "" {
				nickname = mobile
			}
			res, err := session.ExecCtx(ctx,
				"INSERT INTO users (mobile, password, nickname, status, role) VALUES (?,?,?,?,?)",
				mobile, password.Hash(plainPwd), nickname, 1, "merchant_owner",
			)
			if err != nil {
				return err
			}
			ownerID, err = lastInsertID(res)
			if err != nil {
				return err
			}
		} else {
			return err
		}

		shop.OwnerUserID = ownerID
		shop.Status = model.ShopApproved
		res, err := session.ExecCtx(ctx,
			`INSERT INTO shops (name, logo, contact_name, contact_phone, description, category, province, city, district, address, business_license_no, legal_person, license_image, storefront_image, owner_user_id, status)
			 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			shop.Name, shop.Logo, shop.ContactName, shop.ContactPhone, shop.Description, shop.Category,
			shop.Province, shop.City, shop.District, shop.Address, shop.BusinessLicenseNo, shop.LegalPerson,
			shop.LicenseImage, shop.StorefrontImage, shop.OwnerUserID, shop.Status,
		)
		if err != nil {
			return err
		}
		shopID, err := lastInsertID(res)
		if err != nil {
			return err
		}
		shop.ID = shopID
		_, err = session.ExecCtx(ctx,
			"INSERT INTO shop_members (shop_id, user_id, member_role) VALUES (?,?,?)",
			shop.ID, ownerID, model.MemberOwner,
		)
		return err
	})
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (r *MerchantRepository) ResetOwnerPassword(ctx context.Context, shopID uint64, plainPwd string) error {
	shop, err := r.FindShop(ctx, shopID)
	if err != nil {
		return err
	}
	if shop.OwnerUserID == 0 {
		return errors.New("店铺无店主")
	}
	_, err = r.conn.ExecCtx(ctx,
		"UPDATE users SET password=? WHERE id=?", password.Hash(plainPwd), shop.OwnerUserID,
	)
	return err
}

func (r *MerchantRepository) ListShopsByUser(ctx context.Context, userID uint64) ([]model.Shop, error) {
	var shops []model.Shop
	err := r.conn.QueryRowsCtx(ctx, &shops,
		"SELECT "+shopColumns+" FROM shops s JOIN shop_members sm ON sm.shop_id=s.id WHERE sm.user_id=?",
		userID,
	)
	return shops, err
}

func (r *MerchantRepository) IsShopMember(ctx context.Context, shopID, userID uint64) bool {
	n, err := countQuery(ctx, r.conn,
		"SELECT COUNT(*) FROM shop_members WHERE shop_id=? AND user_id=?", shopID, userID,
	)
	return err == nil && n > 0
}
