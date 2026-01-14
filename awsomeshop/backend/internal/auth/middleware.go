package auth

import (
	"awsomeshop/backend/pkg/config"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Claims JWT 声明结构
type Claims struct {
	UserID     uint   `json:"user_id"`
	EmployeeID string `json:"employee_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, employeeID string, role string) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	// 设置过期时间为 24 小时
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 Claims
	claims := &Claims{
		UserID:     userID,
		EmployeeID: employeeID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 Token
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 Token 是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AuthMiddleware JWT 认证中间件
// 验证用户是否已登录
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 中获取 Token
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未登录",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 验证 Token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token 无效或已过期",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("employee_id", claims.EmployeeID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
// 必须在 AuthMiddleware 之后使用
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未登录",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "无权限访问",
				"code":  "FORBIDDEN",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetTokenCookie 设置 Token Cookie
func SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"token",           // name
		token,             // value
		86400,             // maxAge (24 hours in seconds)
		"/",               // path
		"",                // domain (empty for current domain)
		false,             // secure (false for localhost)
		true,              // httpOnly (防止 XSS 攻击)
	)
	// Note: SameSite=Lax is the default in modern browsers
}

// ClearTokenCookie 清除 Token Cookie
func ClearTokenCookie(c *gin.Context) {
	c.SetCookie(
		"token",  // name
		"",       // value
		-1,       // maxAge (negative to delete)
		"/",      // path
		"",       // domain
		false,    // secure
		true,     // httpOnly
	)
}
