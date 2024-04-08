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
	path := lib.GetConfPath("base")
	fmt.Println("path", path)
	err := lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer lib.Destroy()
	fmt.Println("111")
	router.HttpServerRun()
	fmt.Println("222")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
