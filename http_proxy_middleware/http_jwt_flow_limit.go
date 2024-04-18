package http_proxy_middleware

import (
	"fmt"
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 基于请求信息
func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		if appInfo.Qps > 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("service flow limit %v", c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
