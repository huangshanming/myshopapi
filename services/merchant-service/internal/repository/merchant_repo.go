package repository

import (
	"context"
	"errors"
	"strings"

	"mymall/common/password"
	"mymall/services/merchant-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	shopColumns = "id, name, IFNULL(logo,'') AS logo, IFNULL(contact_name,'') AS contact_name, IFNULL(contact_phone,'') AS contact_phone, IFNULL(description,'') AS description, IFNULL(category,'') AS category, IFNULL(province,'') AS province, IFNULL(city,'') AS city, IFNULL(district,'') AS district, IFNULL(address,'') AS address, IFNULL(latitude,0) AS latitude, IFNULL(longitude,0) AS longitude, IFNULL(local_enabled,0) AS local_enabled, IFNULL(business_license_no,'') AS business_license_no, IFNULL(legal_person,'') AS legal_person, IFNULL(license_image,'') AS license_image, IFNULL(storefront_image,'') AS storefront_image, owner_user_id, status, IFNULL(reject_reason,'') AS reject_reason, created_at, updated_at"
	// shopColumnsS is shopColumns with table alias s. for JOIN queries.
	shopColumnsS = "s.id, s.name, IFNULL(s.logo,'') AS logo, IFNULL(s.contact_name,'') AS contact_name, IFNULL(s.contact_phone,'') AS contact_phone, IFNULL(s.description,'') AS description, IFNULL(s.category,'') AS category, IFNULL(s.province,'') AS province, IFNULL(s.city,'') AS city, IFNULL(s.district,'') AS district, IFNULL(s.address,'') AS address, IFNULL(s.latitude,0) AS latitude, IFNULL(s.longitude,0) AS longitude, IFNULL(s.local_enabled,0) AS local_enabled, IFNULL(s.business_license_no,'') AS business_license_no, IFNULL(s.legal_person,'') AS legal_person, IFNULL(s.license_image,'') AS license_image, IFNULL(s.storefront_image,'') AS storefront_image, s.owner_user_id, s.status, IFNULL(s.reject_reason,'') AS reject_reason, s.created_at, s.updated_at"
	shopAppColumns = "id, user_id, shop_name, IFNULL(contact_name,'') AS contact_name, IFNULL(contact_phone,'') AS contact_phone, IFNULL(description,'') AS description, IFNULL(category,'') AS category, IFNULL(province,'') AS province, IFNULL(city,'') AS city, IFNULL(district,'') AS district, IFNULL(address,'') AS address, IFNULL(business_license_no,'') AS business_license_no, IFNULL(legal_person,'') AS legal_person, IFNULL(license_image,'') AS license_image, IFNULL(storefront_image,'') AS storefront_image, status, IFNULL(reject_reason,'') AS reject_reason, IFNULL(reviewed_by,0) AS reviewed_by, reviewed_at, IFNULL(shop_id,0) AS shop_id, created_at, updated_at"
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
	err := r.conn.QueryRowPartialCtx(ctx, &app,
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
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+shopAppColumns+" FROM shop_applications WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) FindApplication(ctx context.Context, id uint64) (*model.ShopApplication, error) {
	var app model.ShopApplication
	err := r.conn.QueryRowPartialCtx(ctx, &app,
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
		if err := session.QueryRowPartialCtx(ctx, &app,
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
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+shopColumns+" FROM shops WHERE "+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) ListPublicShops(ctx context.Context, page, pageSize int, city string) ([]model.Shop, int64, error) {
	where := "status=?"
	args := []any{model.ShopApproved}
	if c := strings.TrimSpace(city); c != "" {
		where += " AND city=?"
		args = append(args, c)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM shops WHERE "+where, args...)
	if err != nil {
		return nil, 0, err
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	var list []model.Shop
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+shopColumns+" FROM shops WHERE "+where+" ORDER BY id ASC LIMIT ? OFFSET ?",
		listArgs...,
	)
	return list, total, err
}

func (r *MerchantRepository) FindShop(ctx context.Context, id uint64) (*model.Shop, error) {
	var shop model.Shop
	err := r.conn.QueryRowPartialCtx(ctx, &shop,
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
		`UPDATE shops SET name=?, logo=?, contact_name=?, contact_phone=?, description=?, category=?, province=?, city=?, district=?, address=?, latitude=?, longitude=?, local_enabled=?, business_license_no=?, legal_person=?, license_image=?, storefront_image=? WHERE id=?`,
		shop.Name, shop.Logo, shop.ContactName, shop.ContactPhone, shop.Description, shop.Category,
		shop.Province, shop.City, shop.District, shop.Address, nullFloat(shop.Latitude), nullFloat(shop.Longitude), shop.LocalEnabled,
		shop.BusinessLicenseNo, shop.LegalPerson,
		shop.LicenseImage, shop.StorefrontImage, shop.ID,
	)
	return err
}

func (r *MerchantRepository) UpdateShopDisplay(ctx context.Context, shop *model.Shop) error {
	_, err := r.conn.ExecCtx(ctx,
		`UPDATE shops SET name=?, logo=?, contact_name=?, contact_phone=?, description=?, category=?, province=?, city=?, district=?, address=?, latitude=?, longitude=?, local_enabled=?, storefront_image=? WHERE id=?`,
		shop.Name, shop.Logo, shop.ContactName, shop.ContactPhone, shop.Description, shop.Category,
		shop.Province, shop.City, shop.District, shop.Address, nullFloat(shop.Latitude), nullFloat(shop.Longitude), shop.LocalEnabled,
		shop.StorefrontImage, shop.ID,
	)
	return err
}

func nullFloat(v float64) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func (r *MerchantRepository) ListShopImages(ctx context.Context, shopID uint64) ([]model.ShopImage, error) {
	var list []model.ShopImage
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT id, shop_id, url, sort, created_at FROM shop_images WHERE shop_id=? ORDER BY sort ASC, id ASC", shopID,
	)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.ShopImage{}
	}
	return list, nil
}

func (r *MerchantRepository) ReplaceShopImages(ctx context.Context, shopID uint64, urls []string) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, "DELETE FROM shop_images WHERE shop_id=?", shopID); err != nil {
			return err
		}
		for i, u := range urls {
			u = strings.TrimSpace(u)
			if u == "" {
				continue
			}
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO shop_images (shop_id, url, sort) VALUES (?,?,?)", shopID, u, i,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *MerchantRepository) ListLocalShops(ctx context.Context, keyword string, lat, lng float64, sortBy string, page, pageSize int) ([]model.Shop, []float64, int64, error) {
	where := "status=? AND local_enabled=1 AND latitude IS NOT NULL AND longitude IS NOT NULL"
	args := []any{model.ShopApproved}
	if k := strings.TrimSpace(keyword); k != "" {
		where += " AND (name LIKE ? OR address LIKE ? OR category LIKE ? OR city LIKE ?)"
		like := "%" + k + "%"
		args = append(args, like, like, like, like)
	}
	total, err := countQuery(ctx, r.conn, "SELECT COUNT(*) FROM shops WHERE "+where, args...)
	if err != nil {
		return nil, nil, 0, err
	}

	hasUserLoc := lat != 0 || lng != 0
	distExpr := "0"
	var selectArgs []any
	if hasUserLoc {
		distExpr = `(6371 * ACOS(LEAST(1, GREATEST(-1,
			COS(RADIANS(?)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(?))
			+ SIN(RADIANS(?)) * SIN(RADIANS(latitude))
		))))`
		selectArgs = append(selectArgs, lat, lng, lat)
	}
	order := "id DESC"
	if hasUserLoc && sortBy == "distance" {
		order = "distance_km ASC, id DESC"
	}

	listSQL := "SELECT " + shopColumns + ", " + distExpr + " AS distance_km FROM shops WHERE " + where +
		" ORDER BY " + order + " LIMIT ? OFFSET ?"
	listArgs := append(selectArgs, args...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)

	type row struct {
		model.Shop
		DistanceKm float64 `db:"distance_km"`
	}
	var rows []row
	if err := r.conn.QueryRowsPartialCtx(ctx, &rows, listSQL, listArgs...); err != nil {
		return nil, nil, 0, err
	}
	list := make([]model.Shop, 0, len(rows))
	dists := make([]float64, 0, len(rows))
	for _, r0 := range rows {
		list = append(list, r0.Shop)
		dists = append(dists, r0.DistanceKm)
	}
	return list, dists, total, nil
}

func (r *MerchantRepository) CreateShopWithOwner(ctx context.Context, shop *model.Shop, mobile, plainPwd, nickname string) (*model.Shop, error) {
	err := r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var ownerID uint64
		err := session.QueryRowPartialCtx(ctx, &ownerID,
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
	err := r.conn.QueryRowsPartialCtx(ctx, &shops,
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
