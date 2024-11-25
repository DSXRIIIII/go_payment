package app

import (
	"github.com/dsxriiiii/l3x_pay/order/app/command"
	"github.com/dsxriiiii/l3x_pay/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
