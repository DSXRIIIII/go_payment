package app

import (
	"github/dsxriiiii/l3x_pay/order/app/command"
	"github/dsxriiiii/l3x_pay/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandle
	UpdateOrder command.UpdateOrderHandle
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
