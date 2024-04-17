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
func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + appInfo.AppID)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appCounter.AppID)))
			c.Abort()
			return
		}
		//dayServiceCounter, _ := appCounter.GetDayData(time.Now())
		fmt.Printf("serviceCounter qps:%v,dayCount:%v", appCounter.QPS, appCounter.TotalCount)

		c.Next()
	}
}
