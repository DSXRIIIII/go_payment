package main

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/config"
	"github.com/dsxriiiii/l3x_pay/common/discovery"
	"github.com/dsxriiiii/l3x_pay/common/genproto/stockpb"
	"github.com/dsxriiiii/l3x_pay/common/logging"
	"github.com/dsxriiiii/l3x_pay/common/server"
	"github.com/dsxriiiii/l3x_pay/common/tracing"
	"github.com/dsxriiiii/l3x_pay/stock/ports"
	"github.com/dsxriiiii/l3x_pay/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
	config.ViperInit()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application := service.NewApplication(ctx)
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()
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
