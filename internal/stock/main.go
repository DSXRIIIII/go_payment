package main

import (
	"github.com/dsxriiiii/l3x_pay/common/genproto/stockpb"
	"github.com/dsxriiiii/l3x_pay/common/server"
	"github.com/spf13/viper"
	"github/dsxriiiii/l3x_pay/stock/ports"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			service := ports.NewGrpcServer()
			stockpb.RegisterStockServiceServer(server, service)
		})
	case "http":
	default:
		panic("unexpect server type")
	}

}
