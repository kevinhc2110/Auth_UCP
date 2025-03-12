package main

import (
	"log"

	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/configs"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/db"
)

func main() {
	dsn := configs.GetEnv("DATABASE_URL") 
	database, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	defer database.Close() // Cerrar conexión al salir

	// Aquí sigues con la inicialización del servidor, rutas, etc.
}
