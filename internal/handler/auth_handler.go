package handler

import (
	"net/http"

	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful", "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")
	// Idealnya fetch full profil user dari database via AuthService
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile retrieved",
		"user": gin.H{
			"id":   userId,
			"role": role,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Pada implementasi stateless JWT, logout biasanya dieksekusi dengan menghapus token dari sisi client
	// Namun endpoint ini tetap dikembalikan 200 OK
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
