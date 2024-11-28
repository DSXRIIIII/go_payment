package service

import (
	"context"
	grpcClient "github.com/dsxriiiii/l3x_pay/common/client"
	"github.com/dsxriiiii/l3x_pay/common/metrics"
	"github.com/dsxriiiii/l3x_pay/payment/adapters"
	"github.com/dsxriiiii/l3x_pay/payment/app"
	"github.com/dsxriiiii/l3x_pay/payment/app/command"
	domain "github.com/dsxriiiii/l3x_pay/payment/domain"
	"github.com/dsxriiiii/l3x_pay/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	//memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
