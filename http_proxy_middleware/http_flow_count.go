package http_proxy_middleware

import (
	"fmt"
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

// 匹配接入方式 基于请求信息
func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		//统计项 1 全站 2 服务 3 租户
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		totalCounter.Increase()
		dayCount, _ := totalCounter.GetDayData(time.Now())
		fmt.Printf("totalCounter qps:%v,dayCount:%v", totalCounter.QPS, dayCount)

		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()
		dayServiceCounter, _ := serviceCounter.GetDayData(time.Now())
		fmt.Printf("serviceCounter qps:%v,dayCount:%v", totalCounter.QPS, dayServiceCounter)

		c.Next()
	}
}