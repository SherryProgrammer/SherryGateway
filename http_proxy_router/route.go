package http_proxy_router

import (
	"github.com/SherryProgrammer/SherryGateway/controller"
	"github.com/SherryProgrammer/SherryGateway/http_proxy_middleware"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		controller.OAuthRegister(oauth)
	}

	root := router.Group("/")
	root.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(), //接入方式

		http_proxy_middleware.HTTPFlowCountMiddleware(),
		http_proxy_middleware.HTTPFlowLimitMiddleware(), //统计限流

		http_proxy_middleware.HTTPJwtAuthTokenMiddleware(),
		http_proxy_middleware.HTTPJwtAuthTokenMiddleware(), //流量统计
		http_proxy_middleware.HTTPWhiteListMiddleware(),
		http_proxy_middleware.HTTPBlackListMiddleware(), //黑白名单

		http_proxy_middleware.HTTPHeaderTransferMiddleware(),
		http_proxy_middleware.HTTPStringUriMiddleware(), //内容替换
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),

		http_proxy_middleware.HTTPReverseProxyMiddleware()) //反向代理

	//OAuthRegister
	return router
}
