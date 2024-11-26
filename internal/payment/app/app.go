package app

import "github.com/dsxriiiii/l3x_pay/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}
