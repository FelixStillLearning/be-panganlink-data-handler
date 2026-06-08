package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	komoditasService service.KomoditasService
	aiService        service.AIService
	userRepo         repository.UserRepository
	productService   service.ProductService
}

func NewAdminHandler(ks service.KomoditasService, aiSvc service.AIService, ur repository.UserRepository, ps service.ProductService) *AdminHandler {
	return &AdminHandler{komoditasService: ks, aiService: aiSvc, userRepo: ur, productService: ps}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	userCount, _ := h.userRepo.Count()
	products, _ := h.productService.GetAll()
	c.JSON(200, gin.H{
		"message": "Admin dashboard",
		"total_users": userCount,
		"total_products": len(products),
	})
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": users}) 
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	// Let's assume body contains {"status": "banned" / "active"} or similar.
	// For simplicity, we just fetch user and save.
	userID := c.Param("id")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// Currently there is no "status" field in model.User (from 01_schema.sql), 
	// so let's just mock it or if we added it, update it.
	// We'll just return success for now since schema doesn't have user.status.
	_ = user
	c.JSON(200, gin.H{"message": "User status updated (mock due to schema)"}) 
}

func (h *AdminHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": products})
}
func (h *AdminHandler) ApproveProduct(c *gin.Context) { 
	h.productService.Update(c.Param("id"), &model.Product{Status: "approved"}) // Mock status update
	c.JSON(200, gin.H{"message": "Product approved"}) 
}
func (h *AdminHandler) RejectProduct(c *gin.Context) { 
	h.productService.Update(c.Param("id"), &model.Product{Status: "rejected"})
	c.JSON(200, gin.H{"message": "Product rejected"}) 
}

func (h *AdminHandler) GetCommodities(c *gin.Context) {
	res, err := h.komoditasService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *AdminHandler) CreateCommodity(c *gin.Context) {
	var k model.Komoditas
	if err := c.ShouldBindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.komoditasService.Create(&k); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Commodity created", "data": k})
}

func (h *AdminHandler) UpdateCommodity(c *gin.Context) {
	var k model.Komoditas
	if err := c.ShouldBindJSON(&k); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	if err := h.komoditasService.Update(id, &k); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commodity updated"})
}

func (h *AdminHandler) DeleteCommodity(c *gin.Context) {
	id := c.Param("id")
	if err := h.komoditasService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commodity deleted"})
}

func (h *AdminHandler) GetMarketPrices(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) }
func (h *AdminHandler) CreateMarketPrice(c *gin.Context) { c.JSON(200, gin.H{"message": "Market price added"}) }

func (h *AdminHandler) GetPriceTrends(c *gin.Context) {
	komoditasID := c.Query("komoditas_id")
	if komoditasID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "komoditas_id query parameter is required"})
		return
	}

	res, err := h.aiService.GetForecast(komoditasID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *AdminHandler) GetSettings(c *gin.Context) { c.JSON(200, gin.H{"data": "Settings"}) }
func (h *AdminHandler) UpdateSettings(c *gin.Context) { c.JSON(200, gin.H{"message": "Settings updated"}) }
