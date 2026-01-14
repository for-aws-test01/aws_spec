package auth

import (
	"awsomeshop/backend/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func init() {
	// 初始化测试配置
	config.AppConfig = &config.Config{
		JWTSecret: "test-secret-key-for-testing",
	}
	gin.SetMode(gin.TestMode)
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint
		employeeID string
		role       string
		wantErr    bool
	}{
		{
			name:       "生成员工 Token",
			userID:     1,
			employeeID: "20260114-001",
			role:       "employee",
			wantErr:    false,
		},
		{
			name:       "生成管理员 Token",
			userID:     2,
			employeeID: "admin",
			role:       "admin",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.employeeID, tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// 生成有效的 Token
	validToken, err := GenerateToken(1, "20260114-001", "employee")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 生成过期的 Token
	expiredClaims := &Claims{
		UserID:     1,
		EmployeeID: "20260114-001",
		Role:       "employee",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte("test-secret-key-for-testing"))

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有效的 Token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "过期的 Token",
			token:   expiredTokenString,
			wantErr: true,
		},
		{
			name:    "无效的 Token",
			token:   "invalid.token.string",
			wantErr: true,
		},
		{
			name:    "空 Token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claims == nil {
				t.Error("ValidateToken() returned nil claims for valid token")
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	// 生成有效的 Token
	validToken, _ := GenerateToken(1, "20260114-001", "employee")

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "有效的 Token",
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: validToken,
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "缺少 Token",
			setupRequest: func(req *http.Request) {
				// 不设置 Cookie
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name: "无效的 Token",
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "invalid.token.string",
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.GET("/test", AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupRequest(req)

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应状态码
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		employeeID     string
		role           string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "管理员访问",
			userID:         1,
			employeeID:     "admin",
			role:           "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "员工访问（无权限）",
			userID:         2,
			employeeID:     "20260114-001",
			role:           "employee",
			expectedStatus: http.StatusForbidden,
			expectedCode:   "FORBIDDEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 生成 Token
			token, _ := GenerateToken(tt.userID, tt.employeeID, tt.role)

			// 创建测试路由
			router := gin.New()
			router.GET("/admin", AuthMiddleware(), AdminMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应状态码
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSetTokenCookie(t *testing.T) {
	router := gin.New()
	router.GET("/set-cookie", func(c *gin.Context) {
		SetTokenCookie(c, "test-token-value")
		c.JSON(http.StatusOK, gin.H{"message": "cookie set"})
	})

	req := httptest.NewRequest(http.MethodGet, "/set-cookie", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证 Cookie 是否设置
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	cookie := cookies[0]
	if cookie.Name != "token" {
		t.Errorf("Expected cookie name 'token', got '%s'", cookie.Name)
	}
	if cookie.Value != "test-token-value" {
		t.Errorf("Expected cookie value 'test-token-value', got '%s'", cookie.Value)
	}
	if cookie.MaxAge != 86400 {
		t.Errorf("Expected MaxAge 86400, got %d", cookie.MaxAge)
	}
	if !cookie.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}
}

func TestClearTokenCookie(t *testing.T) {
	router := gin.New()
	router.GET("/clear-cookie", func(c *gin.Context) {
		ClearTokenCookie(c)
		c.JSON(http.StatusOK, gin.H{"message": "cookie cleared"})
	})

	req := httptest.NewRequest(http.MethodGet, "/clear-cookie", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证 Cookie 是否被清除
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	cookie := cookies[0]
	if cookie.Name != "token" {
		t.Errorf("Expected cookie name 'token', got '%s'", cookie.Name)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge -1, got %d", cookie.MaxAge)
	}
}
