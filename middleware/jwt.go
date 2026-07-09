package middleware

import (
	"mymall/pkg/config"
	"mymall/pkg/jwt"

	"github.com/gin-gonic/gin"
)

var appJWTConfig jwt.Config

func InitJWT(cfg config.JWTConfig) {
	appJWTConfig = jwt.Config{
		Secret:      cfg.Secret,
		ConsumerKey: cfg.ConsumerKey,
		ExpireHours: cfg.ExpireHours,
		Issuer:      cfg.Issuer,
	}
}

func JWTAuth() gin.HandlerFunc {
	return jwt.AuthMiddleware(appJWTConfig.Secret)
}

func GenerateToken(userID uint64, role string) (string, error) {
	return jwt.GenerateToken(userID, role, appJWTConfig)
}
