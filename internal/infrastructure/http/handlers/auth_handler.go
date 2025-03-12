package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinhc2110/Degree-project-UCP/internal/usecases"
)

type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

// NewAuthHandler crea un nuevo manejador de autenticación
func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}	
}

// Login maneja la autenticación del usuario y genera tokens
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validar la solicitud
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Obtener IP y User-Agent del cliente
	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	// Autenticar usuario
	session, accessToken, err := h.authUseCase.Authenticate(c.Request.Context(), req.Email, req.Password, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Responder con los tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": session.RefreshToken,
	})
}

// RefreshToken renueva el token de acceso
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token requerido"})
		return
	}

	newToken, err := h.authUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newToken})
}

// Logout cierra la sesión del usuario
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token requerido"})
		return
	}

	if err := h.authUseCase.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cerrar sesión"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}
