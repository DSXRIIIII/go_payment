package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github/dsxriiiii/l3x_pay/order/adapters"
	"github/dsxriiiii/l3x_pay/order/app"
	"github/dsxriiiii/l3x_pay/order/app/command"
	"github/dsxriiiii/l3x_pay/order/app/query"
	"github/dsxriiiii/l3x_pay/order/metrics"
)

func NewApplication(ctx context.Context) app.Application {
	orderInmemRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandle(orderInmemRepo, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandle(orderInmemRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderInmemRepo, logger, metricClient),
		},
	}
}
