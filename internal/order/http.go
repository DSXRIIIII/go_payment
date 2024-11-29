package main

import (
	"fmt"
	client "github.com/dsxriiiii/l3x_pay/common/client/order"
	"github.com/dsxriiiii/l3x_pay/common/tracing"
	"github.com/dsxriiiii/l3x_pay/order/app"
	"github.com/dsxriiiii/l3x_pay/order/app/command"
	"github.com/dsxriiiii/l3x_pay/order/app/query"
	"github.com/dsxriiiii/l3x_pay/order/convertor"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpServer struct {
	app app.Application
}

// PostCustomerCustomerIDOrders (POST /customer/{customerID}/orders)
func (s HttpServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
	defer span.End()

	var req client.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	r, err := s.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"trace_id":     tracing.TraceID(ctx),
		"customer_D":   req.CustomerID,
		"order_iD":     r.OrderID,
		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
	})
}

// GetCustomerCustomerIDOrdersOrderID (GET /customer/{customerID}/orders/{orderID})
func (s HttpServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
	defer span.End()
	handle, err := s.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		OrderID:    customerID,
		CustomerID: orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"trace_id": tracing.TraceID(ctx),
		"data": gin.H{
			"Order": handle,
		},
	})
}
