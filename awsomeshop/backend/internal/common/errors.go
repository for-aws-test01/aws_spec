package common

// 错误码常量定义
// 命名规则: 大写字母 + 下划线
const (
	// 用户相关错误码
	ErrUserNotFound         = "USER_NOT_FOUND"         // 用户不存在
	ErrInvalidPassword      = "INVALID_PASSWORD"       // 密码错误
	ErrAccountDisabled      = "ACCOUNT_DISABLED"       // 账号已禁用
	ErrEmployeeLimitReached = "EMPLOYEE_LIMIT_REACHED" // 当天员工创建数量已达上限

	// 产品相关错误码
	ErrProductNotFound = "PRODUCT_NOT_FOUND" // 产品不存在
	ErrProductOffline  = "PRODUCT_OFFLINE"   // 产品已下架

	// 积分相关错误码
	ErrInsufficientPoints = "INSUFFICIENT_POINTS" // 积分不足

	// 订单相关错误码
	ErrOrderNotFound        = "ORDER_NOT_FOUND"        // 订单不存在
	ErrOrderCannotCancel    = "ORDER_CANNOT_CANCEL"    // 订单无法取消
	ErrRejectReasonRequired = "REJECT_REASON_REQUIRED" // 拒绝原因必填

	// 认证授权相关错误码
	ErrUnauthorized = "UNAUTHORIZED" // 未授权
	ErrForbidden    = "FORBIDDEN"    // 无权限

	// 通用错误码
	ErrInvalidInput  = "INVALID_INPUT"  // 输入参数错误
	ErrInternalError = "INTERNAL_ERROR" // 内部服务器错误

	// 文件相关错误码
	ErrFileTooLarge = "FILE_TOO_LARGE" // 文件过大
	ErrFileNotFound = "FILE_NOT_FOUND" // 文件不存在
)
