package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/kevinhc2110/Auth_UCP/internal/domain"
	"github.com/kevinhc2110/Auth_UCP/internal/infrastructure/security"
	"github.com/kevinhc2110/Auth_UCP/internal/repositories"
)


type AuthUseCase struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

// NewAuthUseCase crea una nueva instancia del caso de uso de autenticación
func NewAuthUseCase(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// Authenticate valida las credenciales del usuario y genera tokens
func (uc *AuthUseCase) Authenticate(ctx context.Context, email, password, userAgent, clientIP string) (*domain.Session, string, error) {
	// Buscar usuario por email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Verificar contraseña
	if !security.ComparePassword(user.Password, password) {
		return nil, "", ErrInvalidCredentials
	}

	// Generar token JWT
	accessToken, err := security.GenerateToken(user.ID.String(), user.Role, 12*time.Hour)
	if err != nil {
		return nil, "", errors.New("error generating access token")
	}

	// Generar refresh token
	refreshToken := security.GenerateRefreshToken()

	// Crear sesión
	session := &domain.Session{
		ID:           user.ID.String(),
		UserID:       user.ID.String(),
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(12 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Guardar sesión en el repositorio
	err = uc.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, "", errors.New("error saving session")
	}

	return session, accessToken, nil
}

// RefreshToken permite renovar el token de acceso con un refresh token
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	session, err := uc.sessionRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return "", ErrSessionNotFound
	}

	// Validar sesión
	if session.IsBlocked {
		return "", errors.New("session is blocked")
	}
	if time.Now().After(session.ExpiresAt) {
		return "", ErrSessionExpired
	}

	// Generar nuevo token de acceso
	accessToken, err := security.GenerateToken(session.UserID, "user", 12*time.Hour)
	if err != nil {
		return "", errors.New("error generating new access token")
	}

	// Actualizar la sesión con una nueva fecha de expiración
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	session.UpdatedAt = time.Now()

	err = uc.sessionRepo.UpdateSession(ctx, session)
	if err != nil {
		return "", errors.New("error updating session")
	}

	return accessToken, nil
}

// Logout elimina la sesión del usuario
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	return uc.sessionRepo.DeleteSession(ctx, refreshToken)
}
