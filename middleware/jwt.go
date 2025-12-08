package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)
import jwt "github.com/golang-jwt/jwt/v5"

const mySigningKey = "AllYourBase"

// MyCustomClaims 自定义Claims（参考示例中的定义，嵌入RegisteredClaims）
type MyCustomClaims struct {
	Foo string `json:"foo"` // 自定义字段，示例中要求值为"bar"
	jwt.RegisteredClaims
}

// Validate 自定义校验逻辑（参考示例中的Validate方法）
func (m MyCustomClaims) Validate() error {
	// 示例：强制要求Foo字段值为"bar"，可根据业务修改
	if m.Foo != "bar" {
		return errors.New("custom claim validation failed: foo must be 'bar'")
	}
	return nil
}

// -------------------------- 2. 生成Token工具函数（用于测试登录） --------------------------
// GenerateToken 生成JWT Token（参考示例中的NewWithClaims用法）
func GenerateToken() (string, error) {
	// 构造自定义Claims（参考示例中的配置）
	claims := MyCustomClaims{
		Foo: "bar", // 必须为"bar"，否则自定义校验失败
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",                    // 签发者（需与校验逻辑一致时可校验）
			Subject:   "somebody",                // 主题
			ID:        "1",                       // Token ID
			Audience:  []string{"somebody_else"}, // 受众
		},
	}

	// 生成Token（HS256算法，参考示例）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// -------------------------- 3. Gin JWT校验中间件（核心） --------------------------
// JWTAuth Gin中间件：校验JWT Token（完全对齐示例中的解析/校验逻辑）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 步骤1：从请求头提取Token（格式：Bearer <token>）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请求头缺少Authorization字段",
			})
			c.Abort() // 终止请求链，不再执行后续处理
			return
		}

		// 拆分Bearer和Token（参考示例的Token格式）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization格式错误（正确格式：Bearer <token>）",
			})
			c.Abort()
			return
		}
		tokenString := parts[1] // 提取纯Token字符串

		// 步骤2：解析Token（参考示例中的ParseWithClaims用法）
		token, err := jwt.ParseWithClaims(
			tokenString,
			&MyCustomClaims{}, // 绑定自定义Claims
			func(token *jwt.Token) (interface{}, error) {
				// 校验签名算法（必须与生成时一致，防止算法伪造）
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(mySigningKey), nil // 返回签名密钥
			},
			jwt.WithLeeway(5*time.Second), // 时间容错（参考示例中的validationOptions）
		)

		// 步骤3：处理解析错误（参考示例中的errorChecking用法）
		if err != nil {
			// 细分错误类型，返回友好提示
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "无效的Token格式（非标准JWT）"})
			case errors.Is(err, jwt.ErrTokenExpired):
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token已过期"})
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token暂未生效"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token校验失败：" + err.Error()})
			}
			c.Abort()
			return
		}

		// 步骤4：校验Token有效性 + 自定义Claims校验（参考示例中的customValidation）
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			// 执行自定义校验逻辑（示例中的Validate方法）
			if err := claims.Validate(); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "自定义Claims校验失败：" + err.Error()})
				c.Abort()
				return
			}

			// 步骤5：将解析后的Claims存入Gin上下文，供后续接口使用
			c.Set("jwt_claims", claims)
			c.Set("foo", claims.Foo)
			c.Set("issuer", claims.Issuer)
			c.Set("token_id", claims.ID)

			// 校验通过，放行请求
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token无效或签名错误"})
			c.Abort()
		}
	}
}
