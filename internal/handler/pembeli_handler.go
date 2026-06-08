package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type PembeliHandler struct {
	orderService service.OrderService
}

func NewPembeliHandler(os service.OrderService) *PembeliHandler {
	return &PembeliHandler{orderService: os}
}

func (h *PembeliHandler) Dashboard(c *gin.Context) { c.JSON(200, gin.H{"message": "Pembeli dashboard"}) }

func (h *PembeliHandler) GetProducts(c *gin.Context) { c.JSON(200, gin.H{"message": "List of products to buy"}) }

type CheckoutRequest struct {
	Items []model.OrderItem `json:"items" binding:"required"`
}

func (h *PembeliHandler) Checkout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buyerID := c.GetString("user_id")

	order, err := h.orderService.Checkout(buyerID, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    order,
	})
}

func (h *PembeliHandler) GetOrders(c *gin.Context) {
	buyerID := c.GetString("user_id")
	orders, err := h.orderService.GetBuyerOrders(buyerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *PembeliHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.UpdateOrderStatus(orderID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func (h *PembeliHandler) GetHistory(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) }
func (h *PembeliHandler) GetProfile(c *gin.Context) { c.JSON(200, gin.H{"data": "Profile info"}) }
func (h *PembeliHandler) UpdateProfile(c *gin.Context) { c.JSON(200, gin.H{"message": "Profile updated"}) }
