package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/http/handlers"
)

// SetupRoutes define las rutas de la API
func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) {
	api := router.Group("/api")
	{
		// Rutas de autenticación
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)
		api.POST("/logout", authHandler.Logout)

		// Rutas de usuario
		api.POST("/register", userHandler.CreateUser)
		// api.PUT("/users", userHandler.UpdateUser)
		// api.DELETE("/users/:id", userHandler.DeleteUser)
	}
}