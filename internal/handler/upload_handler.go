package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/example/be-panganlink-data-handler/pkg/storage"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	cloudinaryHelper *storage.CloudinaryHelper
}

func NewUploadHandler(ch *storage.CloudinaryHelper) *UploadHandler {
	return &UploadHandler{cloudinaryHelper: ch}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	if h.cloudinaryHelper == nil {
		// Ensure directory exists
		os.MkdirAll("public/uploads", os.ModePerm)
		
		// Fallback to local storage if Cloudinary is not configured
		err = c.SaveUploadedFile(header, "public/uploads/"+header.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally: " + err.Error()})
			return
		}
		
		// Return dynamic local URL based on host
		scheme := "http"
		if c.Request.Header.Get("X-Forwarded-Proto") != "" {
			scheme = c.Request.Header.Get("X-Forwarded-Proto")
		} else if c.Request.TLS != nil {
			scheme = "https"
		}
		host := c.Request.Host
		fileURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, header.Filename)
		c.JSON(http.StatusOK, gin.H{"message": "File uploaded locally", "url": fileURL})
		return
	}

	fileURL, err := h.cloudinaryHelper.UploadFile(c.Request.Context(), file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "url": fileURL})
}
