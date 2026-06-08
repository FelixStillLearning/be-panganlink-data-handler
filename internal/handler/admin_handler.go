package handler

import (
	"net/http"
	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	komoditasService service.KomoditasService
}

func NewAdminHandler(ks service.KomoditasService) *AdminHandler {
	return &AdminHandler{komoditasService: ks}
}

func (h *AdminHandler) Dashboard(c *gin.Context) { c.JSON(200, gin.H{"message": "Admin dashboard data"}) }
func (h *AdminHandler) GetUsers(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) }
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) { c.JSON(200, gin.H{"message": "User status updated"}) }

func (h *AdminHandler) GetProducts(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) }
func (h *AdminHandler) ApproveProduct(c *gin.Context) { c.JSON(200, gin.H{"message": "Product approved"}) }
func (h *AdminHandler) RejectProduct(c *gin.Context) { c.JSON(200, gin.H{"message": "Product rejected"}) }

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

func (h *AdminHandler) GetPriceTrends(c *gin.Context) { c.JSON(200, gin.H{"data": "Price trends from AI"}) }

func (h *AdminHandler) GetSettings(c *gin.Context) { c.JSON(200, gin.H{"data": "Settings"}) }
func (h *AdminHandler) UpdateSettings(c *gin.Context) { c.JSON(200, gin.H{"message": "Settings updated"}) }
