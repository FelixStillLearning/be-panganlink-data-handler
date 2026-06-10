package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	orderService service.OrderService
}

func NewPaymentHandler(os service.OrderService) *PaymentHandler {
	return &PaymentHandler{orderService: os}
}

func (h *PaymentHandler) Webhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.orderService.HandleMidtransWebhook(payload)
	if err != nil {
		if err.Error() == "invalid signature" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
