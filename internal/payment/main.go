package main

import (
	"context"
	"github.com/dsxriiiii/l3x_pay/common/broker"
	"github.com/dsxriiiii/l3x_pay/common/config"
	"github.com/dsxriiiii/l3x_pay/common/logging"
	"github.com/dsxriiiii/l3x_pay/common/server"
	"github.com/dsxriiiii/l3x_pay/common/tracing"
	"github.com/dsxriiiii/l3x_pay/payment/infrastructure/consumer"
	"github.com/dsxriiiii/l3x_pay/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("payment.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverType := viper.GetString("payment.server-to-run")

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)

	switch serverType {
	case "http":
		server.RunHttpServer(serviceName, paymentHandler.RegisterRouters)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}
