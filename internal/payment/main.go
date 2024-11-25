package main

import (
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
