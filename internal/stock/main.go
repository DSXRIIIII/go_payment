package main

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/config"
	"github.com/dsxriiiii/l3x_pay/common/genproto/stockpb"
	"github.com/dsxriiiii/l3x_pay/common/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github/dsxriiiii/l3x_pay/stock/ports"
	"github/dsxriiiii/l3x_pay/stock/service"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := service.NewApplication(ctx)
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
	default:
		panic("unexpect server type")
	}

}
