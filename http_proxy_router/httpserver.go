package http_proxy_router

import (
	"context"
	"github.com/SherryProgrammer/SherryGateway/cert_file"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/SherryProgrammer/go_evnconfig/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}

	log.Printf(" [INFO] http_proxy_run:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil {
		log.Fatalf(" [ERROR] http_proxy_run:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpsServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] https_proxy_run:%s\n", lib.GetStringConf("proxy.https.addr"))
	if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil {
		log.Fatalf(" [ERROR] https_proxy_run:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop stopped\n")
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop stopped\n")
}
