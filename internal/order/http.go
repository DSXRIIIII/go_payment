package main

import "github.com/gin-gonic/gin"

type HttpServer struct{}

// PostCustomerCustomerIDOrders (POST /customer/{customerID}/orders)
func (s HttpServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	// TODO
}

// GetCustomerCustomerIDOrdersOrderID (GET /customer/{customerID}/orders/{orderID})
func (s HttpServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	// TODO
}
