package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevinhc2110/Auth_UCP/internal/infrastructure/http/handlers"
)

// Server representa el servidor HTTP
type Server struct {
	router *gin.Engine
}

// NewServer inicializa un nuevo servidor con los handlers correspondientes
func NewServer(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) *Server {
	router := gin.Default()

	// Registrar rutas con los handlers
	SetupRoutes(router, authHandler, userHandler)

	return &Server{router: router}
}

// Run inicia el servidor en el puerto especificado
func (s *Server) Run(port string) {
	log.Printf("Servidor corriendo en http://localhost:%s", port)
	if err := s.router.Run(":" + port); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
