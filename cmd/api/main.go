package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/configs"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/db"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/http"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/http/handlers"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/security"
	"github.com/kevinhc2110/Degree-project-UCP/internal/usecases"
)

func main() {

	// Cargar claves RSA al iniciar la aplicación
	err := security.LoadKeys()
	if err != nil {
		log.Fatalf("Error cargando las claves RSA: %v", err)
	}

	// Cargar variables de entorno
	configs.LoadEnv()

	// Conectar a la base de datos
	dsn := configs.GetEnv("DATABASE_URL", "DATABASE_URL")
	database, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer database.Close() // Cerrar conexión al salir

	// Crear repositorio de usuarios
	userRepo := db.NewUserRepositoryPg(database)
	sessionRepo := db.NewSessionRepositorypg(database)

	// Crear caso de uso de usuario
	userUseCase := usecases.NewUserUseCase(userRepo)
	authUseCase := usecases.NewAuthUseCase(userRepo, sessionRepo)

	// Crear handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Crear servidor y configurar rutas
	router := gin.Default()
	http.SetupRoutes(router, authHandler, userHandler)

	// Ejecutar el servidor en el puerto 8080
	server := http.NewServer(authHandler, userHandler)

	port := configs.GetEnv("PORT", "8080") // Usa 8080 si no está en .env
	server.Run(port)
}
