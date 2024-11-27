package main

import (
	"fmt"
	"github.com/dsxriiiii/l3x_pay/common/genproto/orderpb"
	"github.com/dsxriiiii/l3x_pay/order/app"
	"github.com/dsxriiiii/l3x_pay/order/app/command"
	"github.com/dsxriiiii/l3x_pay/order/app/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpServer struct {
	app app.Application
}

// PostCustomerCustomerIDOrders (POST /customer/{customerID}/orders)
func (s HttpServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	r, err := s.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"customer_D":   req.CustomerID,
		"order_iD":     r.OrderID,
		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
	})
}

// GetCustomerCustomerIDOrdersOrderID (GET /customer/{customerID}/orders/{orderID})
func (s HttpServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	handle, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    customerID,
		CustomerID: orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"Order": handle,
		},
	})
}
