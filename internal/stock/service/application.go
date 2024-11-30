package service

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/metrics"
	"github.com/dsxriiiii/l3x_pay/stock/adapters"
	"github.com/dsxriiiii/l3x_pay/stock/app"
	"github.com/dsxriiiii/l3x_pay/stock/app/query"
	"github.com/dsxriiiii/l3x_pay/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	stripeAPI := integration.NewStripeAPI()
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
