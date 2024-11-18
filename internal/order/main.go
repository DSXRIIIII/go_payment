package main

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/config"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
	"github.com/dsxriiiii/l3x_pay/common/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github/dsxriiiii/l3x_pay/order/ports"
	"github/dsxriiiii/l3x_pay/order/service"
	"google.golang.org/grpc"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HttpServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
