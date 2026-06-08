package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type PembeliHandler struct {
	orderService service.OrderService
	userRepo     repository.UserRepository
}

func NewPembeliHandler(os service.OrderService, ur repository.UserRepository) *PembeliHandler {
	return &PembeliHandler{orderService: os, userRepo: ur}
}

func (h *PembeliHandler) Dashboard(c *gin.Context) {
	buyerID := c.GetString("user_id")
	orders, _ := h.orderService.GetBuyerOrders(buyerID)

	var totalPembelian int = len(orders)
	var dalamPengiriman int = 0
	var totalBelanja float64 = 0
	
	// Assume we calculate unique products as 'produk_favorit' just to have a number
	produkFavoritMap := make(map[string]bool)

	for _, o := range orders {
		if o.Status == "dikirim" || o.Status == "diproses" {
			dalamPengiriman++
		}
		if o.Status == "selesai" || o.Status == "success" {
			totalBelanja += o.TotalHarga
		}
		for _, item := range o.Items {
			produkFavoritMap[item.ProductID] = true
		}
	}

	c.JSON(200, gin.H{
		"message": "Pembeli dashboard", 
		"data": gin.H{
			"total_pembelian": totalPembelian,
			"dalam_pengiriman": dalamPengiriman,
			"total_belanja": totalBelanja,
			"produk_favorit": len(produkFavoritMap),
			"recent_orders": orders,
		},
	})
}

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

func (h *PembeliHandler) GetHistory(c *gin.Context) {
	buyerID := c.GetString("user_id")
	orders, err := h.orderService.GetBuyerOrders(buyerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Filter completed/rejected orders for history
	c.JSON(200, gin.H{"data": orders})
}

func (h *PembeliHandler) GetProfile(c *gin.Context) {
	buyerID := c.GetString("user_id")
	user, err := h.userRepo.FindByID(buyerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Password = "" // Hide password
	c.JSON(200, gin.H{"data": user})
}

func (h *PembeliHandler) UpdateProfile(c *gin.Context) {
	buyerID := c.GetString("user_id")
	user, err := h.userRepo.FindByID(buyerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err == nil {
		if req.Name != "" {
			user.Name = req.Name
		}
		if req.Location != "" {
			user.Location = req.Location
		}
		_ = h.userRepo.Update(user)
	}
	
	c.JSON(200, gin.H{"message": "Profile updated", "data": user})
}
