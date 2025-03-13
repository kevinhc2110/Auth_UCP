package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/http/handlers"
)

// SetupRoutes define las rutas de la API
func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) {
	api := router.Group("/api")

	{
		// Rutas de autenticación y usuarios
		api.POST("/login", authHandler.Login)
		api.POST("/register", userHandler.CreateUser)

	}

	// Rutas protegidas
	protected := api.Group("/")
	protected.Use(AuthMiddleware()) // 🔐 Middleware aplicado

	{
		protected.PUT("/users", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		protected.POST("/refresh", authHandler.RefreshToken)
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/public-key", handlers.PublicKeyHandler)
	}
}
