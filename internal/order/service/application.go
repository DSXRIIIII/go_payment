package service

import (
	"context"
	grpcClient "github.com/dsxriiiii/l3x_pay/common/client"
	"github.com/sirupsen/logrus"
	"github/dsxriiiii/l3x_pay/order/adapters"
	"github/dsxriiiii/l3x_pay/order/adapters/grpc"
	"github/dsxriiiii/l3x_pay/order/app"
	"github/dsxriiiii/l3x_pay/order/app/command"
	"github/dsxriiiii/l3x_pay/order/app/query"
	"github/dsxriiiii/l3x_pay/order/metrics"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
