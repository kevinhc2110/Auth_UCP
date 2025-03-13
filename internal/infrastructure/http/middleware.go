package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/security"
)

// AuthMiddleware verifica el JWT antes de permitir acceso al handler
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// Validar que el token tenga el prefijo "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Extraer el token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := security.ValidateToken(tokenString)
		if err != nil {
			if err.Error() == "token expirado" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			}
			c.Abort()
			return
		}

		// Validar el rol
		if claims.Role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Rol no válido"})
			c.Abort()
			return
		}

		// Guardar los claims en el contexto para usarlos en el handler
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

