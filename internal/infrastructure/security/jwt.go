package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/configs"
)

// Claims personalizados para el token
 type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken genera un JWT para un usuario
func GenerateToken(userID, role string, duration time.Duration) (string, error) {
	// Obtener la clave secreta en tiempo de ejecución
	jwtSecret := []byte(configs.GetEnv("JWT_SECRET_KEY"))
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret is not set")
	}

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // Expiración
			IssuedAt:  jwt.NewNumericDate(time.Now()),              // Fecha de emisión
		},
	}

	// Crear el token con los claims y firmarlo
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken verifica la validez de un JWT
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Obtener la clave secreta en tiempo de ejecución
	jwtSecret := []byte(configs.GetEnv("JWT_SECRET_KEY"))
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma no válido")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extraer los claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Verificar si ha expirado
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expirado")
	}

	return claims, nil
}

// GenerateRefreshToken genera un token de refresco único
func GenerateRefreshToken() string {
	return uuid.New().String()
}
