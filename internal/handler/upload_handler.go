package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/pkg/storage"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	azureHelper *storage.AzureHelper
}

func NewUploadHandler(ah *storage.AzureHelper) *UploadHandler {
	return &UploadHandler{azureHelper: ah}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	if h.azureHelper == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Azure storage is not configured"})
		return
	}

	contentType := header.Header.Get("Content-Type")
	fileURL, err := h.azureHelper.UploadFile(c.Request.Context(), file, header.Filename, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "url": fileURL})
}
