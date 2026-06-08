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
	testimonials := []map[string]string{
		{
			"text": "Semenjak menggunakan PanganLink, saya tahu pasti harga pasaran jagung saya lewat fitur prediksi AI. Tidak ada lagi tengkulak yang mempermainkan harga panen saya.",
			"name": "Pak Wahyudi",
			"role": "Petani Jagung, Jawa Timur",
			"initial": "W",
		},
		{
			"text": "Sangat membantu bisnis katering kami! Harga bahan pokok jauh lebih stabil dan wajar dibanding beli dari pasar tangan ketiga. Kualitas sayurnya pun dijamin segar.",
			"name": "Ibu Ningsih",
			"role": "Pemilik Katering, Jakarta",
			"initial": "I",
		},
	}
	c.JSON(200, gin.H{"data": testimonials})
}
