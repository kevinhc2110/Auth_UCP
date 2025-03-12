package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims personalizados para el token
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken genera un JWT para un usuario
func GenerateToken(userID, role string, duration time.Duration) (string, error) {
	if privateKey == nil {
		return "", errors.New("clave privada no cargada")
	}

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ValidateToken verifica la validez de un JWT
func ValidateToken(tokenString string) (*JWTClaims, error) {
	if publicKey == nil {
		return nil, errors.New("clave pública no cargada")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("método de firma no válido")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extrael claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Verificar expiración
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expirado")
	}

	return claims, nil
}

// GenerateRefreshToken genera un token de refresco único
func GenerateRefreshToken() string {
	return uuid.New().String()
}
