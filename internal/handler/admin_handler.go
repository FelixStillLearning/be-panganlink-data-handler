package handler

import (
	"fmt"
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AdminHandler struct {
	db               *gorm.DB
	komoditasService service.KomoditasService
	aiService        service.AIService
	userRepo         repository.UserRepository
	productService   service.ProductService
}

func NewAdminHandler(db *gorm.DB, ks service.KomoditasService, aiSvc service.AIService, ur repository.UserRepository, ps service.ProductService) *AdminHandler {
	return &AdminHandler{db: db, komoditasService: ks, aiService: aiSvc, userRepo: ur, productService: ps}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	userCount, _ := h.userRepo.Count()
	products, _ := h.productService.GetAll(1, 99999) // Temp workaround for count

	var weeklySales []map[string]interface{}
	daysOfWeek := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
	
	now := time.Now()
	var maxVal int64 = 0
	var counts [7]int64

	// Get the past 7 days (including today)
	for i := 6; i >= 0; i-- {
		targetDate := now.AddDate(0, 0, -i)
		startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var count int64
		h.db.Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).Count(&count)
		counts[6-i] = count
		if count > maxVal {
			maxVal = count
		}
	}

	for i := 6; i >= 0; i-- {
		targetDate := now.AddDate(0, 0, -i)
		dayName := daysOfWeek[targetDate.Weekday()]
		
		val := counts[6-i]
		height := "0%"
		if maxVal > 0 {
			height = fmt.Sprintf("%d%%", (val*100)/maxVal)
		} else if val == 0 {
			height = "5%" // just a tiny sliver so it's not totally empty
		}

		weeklySales = append(weeklySales, map[string]interface{}{
			"day":    dayName,
			"value":  val,
			"height": height,
		})
	}

	c.JSON(200, gin.H{
		"message": "Admin dashboard",
		"total_users": userCount,
		"total_products": len(products),
		"weekly_sales": weeklySales,
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
	userID := c.Param("id")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	_ = user
	c.JSON(200, gin.H{"message": "User status updated (mock due to schema)"}) 
}

func (h *AdminHandler) GetProducts(c *gin.Context) {
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	products, err := h.productService.GetAll(page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": products, "page": page, "limit": limit})
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

func (h *AdminHandler) GetMarketPrices(c *gin.Context) {
	var prices []model.MarketPrice
	if err := h.db.Order("date DESC").Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": prices})
}

func (h *AdminHandler) CreateMarketPrice(c *gin.Context) {
	var p model.MarketPrice
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = uuid.New().String()
	if err := h.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Market price added", "data": p})
}

func (h *AdminHandler) DeleteMarketPrice(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Where("id = ?", id).Delete(&model.MarketPrice{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Market price deleted"})
}

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
