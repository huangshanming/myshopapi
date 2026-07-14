package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

const (
	RoleUser           = "user"
	RolePlatformAdmin  = "platform_admin"
	RoleMerchantOwner  = "merchant_owner"
	RoleMerchantStaff  = "merchant_staff"
	RoleAdmin          = "admin" // legacy alias → treat as platform_admin
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type Claims struct {
	Key    string `json:"key"` // APISIX jwt-auth consumer key
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
	ShopID uint64 `json:"shop_id,omitempty"`
	jwtlib.RegisteredClaims
}

type Config struct {
	Secret      string
	ConsumerKey string
	ExpireHours int
	Issuer      string
}

func GenerateToken(userID uint64, role string, cfg Config) (string, error) {
	return GenerateTokenWithShop(userID, role, 0, cfg)
}

func GenerateTokenWithShop(userID uint64, role string, shopID uint64, cfg Config) (string, error) {
	if cfg.Secret == "" {
		return "", errors.New("jwt secret is required")
	}
	if role == "" {
		role = RoleUser
	}
	if role == RoleAdmin {
		role = RolePlatformAdmin
	}
	expireHours := cfg.ExpireHours
	if expireHours <= 0 {
		expireHours = 24
	}
	issuer := cfg.Issuer
	if issuer == "" {
		issuer = "mymall"
	}
	consumerKey := cfg.ConsumerKey
	if consumerKey == "" {
		consumerKey = "mymall-user-key"
	}

	now := time.Now()
	claims := Claims{
		Key:    consumerKey,
		UserID: userID,
		Role:   role,
		ShopID: shopID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			Issuer:    issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenStr string, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	token, err := jwtlib.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwtlib.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
		jwtlib.WithLeeway(5*time.Second),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwtlib.ErrTokenExpired):
			return nil, ErrExpiredToken
		case errors.Is(err, jwtlib.ErrTokenMalformed),
			errors.Is(err, jwtlib.ErrTokenSignatureInvalid),
			errors.Is(err, jwtlib.ErrTokenNotValidYet):
			return nil, ErrInvalidToken
		default:
			return nil, err
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.UserID == 0 || claims.Key == "" {
		return nil, ErrInvalidToken
	}
	if claims.Role == RoleAdmin {
		claims.Role = RolePlatformAdmin
	}
	return claims, nil
}

func IsPlatformAdmin(role string) bool {
	return role == RolePlatformAdmin || role == RoleAdmin
}

func IsMerchant(role string) bool {
	return role == RoleMerchantOwner || role == RoleMerchantStaff
}
