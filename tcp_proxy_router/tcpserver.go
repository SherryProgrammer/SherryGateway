package tcp_proxy_router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/reverse_proxy"
	"github.com/SherryProgrammer/SherryGateway/tcp_proxy_middleware"
	"github.com/SherryProgrammer/SherryGateway/tcp_server"
	"log"
	"net"
)

var tcpServerList = []*tcp_server.TcpServer{}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TcpServerRun() {

	serviceList := dao.ServiceManagerHandler.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		//log.Printf(" [INFO] tcp_proxy_run:%s\n", tempItem.TCPRule.Port)
		marshal, _ := json.Marshal(tempItem)
		fmt.Println(string(marshal))

		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf("%d", serviceDetail.TCPRule.Port)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}
			//构建路由及设置中间件
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				tcp_proxy_middleware.TCPFlowLimitMiddleware(),
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
				func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
					return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
				}, router)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)

			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf(" [INFO] tcp_proxy_run %s\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf(" [INFO] tcp_proxy_run &v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)

	}
}
