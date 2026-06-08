package handler

import (
	"net/http"
	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type PetaniHandler struct {
	productService service.ProductService
}

func NewPetaniHandler(ps service.ProductService) *PetaniHandler {
	return &PetaniHandler{productService: ps}
}

func (h *PetaniHandler) Dashboard(c *gin.Context) { c.JSON(200, gin.H{"message": "Petani dashboard"}) }

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

func (h *PetaniHandler) GetOrders(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) }
func (h *PetaniHandler) UpdateOrderStatus(c *gin.Context) { c.JSON(200, gin.H{"message": "Order status updated"}) }
func (h *PetaniHandler) GetHistory(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) }

func (h *PetaniHandler) GetRecommendations(c *gin.Context) { c.JSON(200, gin.H{"data": "AI recommendations"}) }

func (h *PetaniHandler) GetProfile(c *gin.Context) { c.JSON(200, gin.H{"data": "Profile info"}) }
func (h *PetaniHandler) UpdateProfile(c *gin.Context) { c.JSON(200, gin.H{"message": "Profile updated"}) }
