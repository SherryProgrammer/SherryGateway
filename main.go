package main

import (
	"flag"
	"fmt"
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/grpc_proxy_router"
	"github.com/SherryProgrammer/SherryGateway/http_proxy_router"
	"github.com/SherryProgrammer/SherryGateway/router"
	"github.com/SherryProgrammer/SherryGateway/tcp_proxy_router"
	"github.com/SherryProgrammer/go_evnconfig/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// endpoint dashboard 后台管理 server 代理服务器
// config ./conf/prod/ 对应配置文件夹
var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("conf", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()

	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		//如果configPath为空 从命令行中`-congig-./conf/prod/`中读取
		//path := lib.GetConfPath("base") // 获取基础模块的配置路径。
		//fmt.Println("path", path)       // 打印配置路径。
		// 使用配置路径和模块名称前缀初始化模块。
		//"./conf/dev/"
		err := lib.InitModule(*config, []string{"base", "mysql", "redis"}) //请求链路日志打印，涵盖mysql/redis/request
		if err != nil {
			log.Fatalf("%v", err) // 如果初始化失败，则退出程序。
		}
		defer lib.Destroy() // 退出
		//fmt.Println("111")
		router.HttpServerRun() // 启动 HTTP 服务器。
		//fmt.Println("222")
		fmt.Println("start server")
		quit := make(chan os.Signal)                                                           // 创建通道以接收操作系统信号。 优雅停止模板
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // ctrl+c kill
		<-quit                                                                                 // 等待接收到信号。
		router.HttpServerStop()                                                                // 停止 HTTP 服务器。
	} else {
		err := lib.InitModule(*config, []string{"base", "mysql", "redis"}) //请求链路日志打印，涵盖mysql/redis/request
		if err != nil {
			log.Fatalf("%v", err) // 如果初始化失败，则退出程序。
		}
		defer lib.Destroy()                  // 退出
		dao.ServiceManagerHandler.LoadOnce() //启动时加载服务列表
		dao.AppManagerHandler.LoadOnce()
		go func() { //携程
			http_proxy_router.HttpServerRun()
		}()
		go func() { //携程
			http_proxy_router.HttpsServerRun()
		}()
		go func() { //携程
			tcp_proxy_router.TcpServerRun()
		}()
		go func() { //携程
			grpc_proxy_router.GrpcServerRun()
		}()
		fmt.Println("start server")
		//todo

		quit := make(chan os.Signal)                                                           // 创建通道以接收操作系统信号。 优雅停止模板
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // ctrl+c kill
		<-quit                                                                                 // 等待接收到信号。

		tcp_proxy_router.TcpServerStop()
		grpc_proxy_router.GrpcServerStop()
		http_proxy_router.HttpServerStop() // 停止 HTTP 服务器。
		http_proxy_router.HttpsServerStop()

	}

}
