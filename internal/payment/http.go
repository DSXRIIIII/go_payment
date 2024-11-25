package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) RegisterRouters(c *gin.Engine) {
	c.POST("api/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	logrus.Info("receive webhook from stripe")
}