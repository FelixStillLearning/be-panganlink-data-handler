package handler

import "github.com/gin-gonic/gin"

func GetPublicStats(c *gin.Context) {
	c.JSON(200, gin.H{"total_petani": 150, "total_transaksi": 1200, "produk_aktif": 340})
}

func GetPublicCommodities(c *gin.Context) {
	c.JSON(200, gin.H{"data": []gin.H{
		{"id": "kmdB-001", "nama": "Beras", "harga": 14500},
		{"id": "kmdC-001", "nama": "Cabai", "harga": 45000},
	}})
}

func GetPublicTestimonials(c *gin.Context) {
	c.JSON(200, gin.H{"data": []string{"Sangat membantu!", "Harga transparan."}})
}
