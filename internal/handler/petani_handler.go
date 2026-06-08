package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type PetaniHandler struct {
	productService service.ProductService
	aiService      service.AIService
	orderService   service.OrderService
	userRepo       repository.UserRepository
}

func NewPetaniHandler(ps service.ProductService, aiSvc service.AIService, os service.OrderService, ur repository.UserRepository) *PetaniHandler {
	return &PetaniHandler{productService: ps, aiService: aiSvc, orderService: os, userRepo: ur}
}

func (h *PetaniHandler) Dashboard(c *gin.Context) {
	petaniID := c.GetString("user_id")
	
	products, _ := h.productService.GetByUserID(petaniID)
	totalProducts := len(products)
	pendingApprove := 0
	for _, p := range products {
		if p.Status == "pending" {
			pendingApprove++
		}
	}

	orders, _ := h.orderService.GetPetaniOrders(petaniID)
	ordersPending := 0
	var salesThisMonth float64 = 0
	
	for _, o := range orders {
		if o.Status == "pending" || o.Status == "menunggu" {
			ordersPending++
		}
		if o.Status == "selesai" || o.Status == "success" {
			salesThisMonth += o.TotalHarga
		}
	}

	c.JSON(200, gin.H{
		"message": "Petani dashboard",
		"data": gin.H{
			"total_products": totalProducts,
			"pending_approve": pendingApprove,
			"orders_pending": ordersPending,
			"sales_this_month": salesThisMonth,
			"recent_orders": orders,
		},
	})
}

func (h *PetaniHandler) GetProducts(c *gin.Context) {
	userId := c.GetString("user_id")
	res, err := h.productService.GetByUserID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *PetaniHandler) CreateProduct(c *gin.Context) {
	var p model.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.UserID = c.GetString("user_id")
	if err := h.productService.Create(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created", "data": p})
}

func (h *PetaniHandler) UpdateProduct(c *gin.Context) {
	var p model.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	p.UserID = c.GetString("user_id")
	if err := h.productService.Update(id, &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func (h *PetaniHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *PetaniHandler) GetOrders(c *gin.Context) {
	petaniID := c.GetString("user_id")
	orders, err := h.orderService.GetPetaniOrders(petaniID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *PetaniHandler) UpdateOrderStatus(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func (h *PetaniHandler) GetHistory(c *gin.Context) { 
	petaniID := c.GetString("user_id")
	orders, _ := h.orderService.GetPetaniOrders(petaniID)
	c.JSON(200, gin.H{"data": orders})
}

func (h *PetaniHandler) GetRecommendations(c *gin.Context) {
	komoditasID := c.Query("komoditas_id")
	if komoditasID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "komoditas_id query parameter is required"})
		return
	}

	res, err := h.aiService.GetRecommendation(komoditasID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *PetaniHandler) GetProfile(c *gin.Context) {
	petaniID := c.GetString("user_id")
	user, err := h.userRepo.FindByID(petaniID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Password = ""
	c.JSON(200, gin.H{"data": user})
}

func (h *PetaniHandler) UpdateProfile(c *gin.Context) {
	petaniID := c.GetString("user_id")
	user, err := h.userRepo.FindByID(petaniID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	// Bind body and update specific fields
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
