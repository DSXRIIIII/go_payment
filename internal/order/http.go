package main

import (
	"fmt"
	"github.com/dsxriiiii/l3x_pay/common"
	client "github.com/dsxriiiii/l3x_pay/common/client/order"
	"github.com/dsxriiiii/l3x_pay/order/app"
	"github.com/dsxriiiii/l3x_pay/order/app/command"
	"github.com/dsxriiiii/l3x_pay/order/app/dto"
	"github.com/dsxriiiii/l3x_pay/order/app/query"
	"github.com/dsxriiiii/l3x_pay/order/convertor"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	common.BaseResponse
	app app.Application
}

// PostCustomerCustomerIdOrders (POST /customer/{customerID}/orders)
func (H HttpServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {

	var (
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
		err  error
	)

	defer func() {
		H.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResponse{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

// GetCustomerCustomerIdOrdersOrderId (GET /customer/{customerID}/orders/{orderID})
func (H HttpServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp interface{}
	)
	defer func() {
		H.Response(c, err, resp)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		return
	}

	resp = convertor.NewOrderConvertor().EntityToClient(o)
}
