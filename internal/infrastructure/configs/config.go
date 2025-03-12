package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables de entorno desde un archivo .env (si existe)
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}
}

// GetEnv obtiene una variable de entorno con una clave específica
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Advertencia: La variable de entorno %s no está definida, usando valor por defecto: %s\n", key, defaultValue)
		return defaultValue
	}
	return value
}
