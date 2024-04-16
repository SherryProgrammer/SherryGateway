package http_proxy_middleware

import (
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
			item := strings.Split(item, "")
			if len(item) != 3 {
				continue
			}
			if item[0] == "add" || item[0] == "edit" {
				c.Request.Header.Set(item[1], item[2])
			}
			if item[0] == "del" {
				c.Request.Header.Del(item[1])
			}
		}
		c.Next()

	}
}
