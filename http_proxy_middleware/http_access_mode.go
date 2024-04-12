package http_proxy_middleware

import (
	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPAccessModelMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

	}
}
