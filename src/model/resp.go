package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Resp 统一响应结构
type Resp struct {
	TraceID string      `json:"trace_id"`       // 事务ID，用于追踪请求
	Code    int         `json:"code"`           // 业务状态码：0-成功，其他-失败
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// 业务状态码定义
const (
	CodeSuccess      = 0    // 成功
	CodeError        = -1   // 通用错误
	CodeParamError   = 1001 // 参数错误
	CodeUnauthorized = 1002 // 未授权
	CodeForbidden    = 1003 // 无权限
	CodeNotFound     = 1004 // 资源不存在
	CodeServerError  = 1005 // 服务器错误
)

// TraceIDKey 存储在 gin.Context 中的 TraceID 键名
const TraceIDKey = "trace_id"

// GenerateTransID 生成事务ID（使用 UUID v4）
func GenerateTransID() string {
	return uuid.New().String()
}

// GetTraceID 从 gin.Context 中获取追踪ID，如果不存在则生成新的
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return GenerateTransID()
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		TraceID: GetTraceID(c),
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		TraceID: GetTraceID(c),
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, businessCode int, message string) {
	c.JSON(httpCode, Resp{
		TraceID: GetTraceID(c),
		Code:    businessCode,
		Message: message,
	})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, httpCode int, businessCode int, message string, data interface{}) {
	c.JSON(httpCode, Resp{
		TraceID: GetTraceID(c),
		Code:    businessCode,
		Message: message,
		Data:    data,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeParamError, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden 无权限响应
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeServerError, message)
}
