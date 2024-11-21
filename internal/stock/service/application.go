package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github/dsxriiiii/l3x_pay/common/metrics"
	"github/dsxriiiii/l3x_pay/stock/adapters"
	"github/dsxriiiii/l3x_pay/stock/app"
	"github/dsxriiiii/l3x_pay/stock/app/query"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
