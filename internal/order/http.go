package main

import (
	"github.com/gin-gonic/gin"
	"github/dsxriiiii/l3x_pay/order/app"
	"github/dsxriiiii/l3x_pay/order/app/query"
	"net/http"
)

type HttpServer struct {
	app app.Application
}

// PostCustomerCustomerIDOrders (POST /customer/{customerID}/orders)
func (s HttpServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	// TODO
}

// GetCustomerCustomerIDOrdersOrderID (GET /customer/{customerID}/orders/{orderID})
func (s HttpServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	handle, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    "fake-ID",
		CustomerID: "fake-customer-id",
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": handle})
}
