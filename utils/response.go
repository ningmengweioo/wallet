package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Created
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "created successfully",
		Data:    data,
	})
}

// Error response
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
