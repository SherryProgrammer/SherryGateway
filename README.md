# SherryGateway
Supports reverse proxy and load balancing for HTTP, TCP, and GRPC

启动控制面服务
go run main.go -conf ./conf/dev/ -endpoint dashboard

启动代理服务器服务
go run main.go -conf ./conf/dev/ -endpoint server
