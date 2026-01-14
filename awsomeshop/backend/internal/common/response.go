package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一成功响应格式
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationData 分页数据结构
type PaginationData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// SuccessResponse 返回成功响应
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Data:    data,
		Message: message,
	})
}

// CreatedResponse 返回创建成功响应 (201)
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, Response{
		Data:    data,
		Message: message,
	})
}

// ErrorResponseWithCode 返回错误响应（自定义状态码）
func ErrorResponseWithCode(c *gin.Context, statusCode int, errorMsg string, errorCode string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Error:   errorMsg,
		Code:    errorCode,
		Details: details,
	})
}

// BadRequestError 返回 400 错误
func BadRequestError(c *gin.Context, errorMsg string, errorCode string) {
	ErrorResponseWithCode(c, http.StatusBadRequest, errorMsg, errorCode, nil)
}

// UnauthorizedError 返回 401 错误
func UnauthorizedError(c *gin.Context, errorMsg string, errorCode string) {
	ErrorResponseWithCode(c, http.StatusUnauthorized, errorMsg, errorCode, nil)
}

// ForbiddenError 返回 403 错误
func ForbiddenError(c *gin.Context, errorMsg string, errorCode string) {
	ErrorResponseWithCode(c, http.StatusForbidden, errorMsg, errorCode, nil)
}

// NotFoundError 返回 404 错误
func NotFoundError(c *gin.Context, errorMsg string, errorCode string) {
	ErrorResponseWithCode(c, http.StatusNotFound, errorMsg, errorCode, nil)
}

// ConflictError 返回 409 错误
func ConflictError(c *gin.Context, errorMsg string, errorCode string) {
	ErrorResponseWithCode(c, http.StatusConflict, errorMsg, errorCode, nil)
}

// InternalServerError 返回 500 错误
func InternalServerError(c *gin.Context, errorMsg string, errorCode string) {
	ErrorResponseWithCode(c, http.StatusInternalServerError, errorMsg, errorCode, nil)
}

// PaginationResponse 返回分页响应
func PaginationResponse(c *gin.Context, items interface{}, total int64, page int, pageSize int, message string) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Data: PaginationData{
			Items:      items,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
		Message: message,
	})
}
