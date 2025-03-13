package security

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// Cargar claves RSA desde archivos
func LoadKeys() error {
	privKeyBytes, err := os.ReadFile("internal/infrastructure/security/private.pem")
	if err != nil {
		return err
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		return err
	}

	pubKeyBytes, err := os.ReadFile("internal/infrastructure/security/public.pem")
	if err != nil {
		return err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	return err
}

// Obtener claves
func PrivateKey() *rsa.PrivateKey {
	return privateKey
}

func PublicKey() *rsa.PublicKey {
	return publicKey
}
