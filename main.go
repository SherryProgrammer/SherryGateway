package main

import (
	"fmt"
	"github.com/SherryProgrammer/SherryGateway/router"
	"github.com/SherryProgrammer/go_evnconfig/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//如果configPath为空 从命令行中`-congig-./conf/prod/`中读取
	path := lib.GetConfPath("base") // 获取基础模块的配置路径。
	fmt.Println("path", path)       // 打印配置路径。
	// 使用配置路径和模块名称初始化模块。
	err := lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"}) //请求链路日志打印，涵盖mysql/redis/request
	if err != nil {
		log.Fatalf("%v", err) // 如果初始化失败，则退出程序。
	}
	defer lib.Destroy() // 延迟销毁已初始化的模块。
	//fmt.Println("111")
	router.HttpServerRun() // 启动 HTTP 服务器。
	//fmt.Println("222")

	quit := make(chan os.Signal)                                                           // 创建通道以接收操作系统信号。
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // 注册信号到 quit 通道。
	<-quit                                                                                 // 等待接收到信号。

	router.HttpServerStop() // 停止 HTTP 服务器。
}
