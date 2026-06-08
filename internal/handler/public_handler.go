package handler

import (
	"net/http"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	komoditasService service.KomoditasService
}

func NewPublicHandler(ks service.KomoditasService) *PublicHandler {
	return &PublicHandler{komoditasService: ks}
}

func (h *PublicHandler) GetStats(c *gin.Context) {
	c.JSON(200, gin.H{"total_petani": 150, "total_transaksi": 1200, "produk_aktif": 340})
}

func (h *PublicHandler) GetCommodities(c *gin.Context) {
	res, err := h.komoditasService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *PublicHandler) GetTestimonials(c *gin.Context) {
	c.JSON(200, gin.H{"data": []string{"Sangat membantu!", "Harga transparan."}})
}
