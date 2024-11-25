package main

import (
	"github.com/dsxriiiii/l3x_pay/common/broker"
	"github.com/dsxriiiii/l3x_pay/common/config"
	"github.com/dsxriiiii/l3x_pay/common/logging"
	"github.com/dsxriiiii/l3x_pay/common/server"
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
	serverType := viper.GetString("payment.server-to-run")

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

	paymentHandler := NewPaymentHandler()

	switch serverType {
	case "http":
		server.RunHttpServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRouters)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}
