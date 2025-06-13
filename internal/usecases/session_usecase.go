package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kevinhc2110/Auth_UCP/internal/domain"
	"github.com/kevinhc2110/Auth_UCP/internal/repositories"
)

var (
	ErrSessionNotFound = errors.New("sessión no encontrada")
	ErrSessionBlocked  = errors.New("sessión esta bloqueada")
	ErrInvalidSession  = errors.New("token de sessión invalido")
	ErrSessionExpired  = errors.New("sessión a expirado")
	ErrInvalidDuration = errors.New("duración de sesión no válida")
	ErrInvalidToken    = errors.New("token invalido")
)

type SessionUseCase struct {
	repo repositories.SessionRepository
}

// NewSessionUseCase crea una nueva instancia del caso de uso de sesión
func NewSessionUseCase(repo repositories.SessionRepository) *SessionUseCase {
	return &SessionUseCase{repo: repo}
}

// CreateSession crea una nueva sesión para un usuario
func (uc *SessionUseCase) CreateSession(ctx context.Context, userID, userAgent, clientIP string, refreshToken string, duration time.Duration) (*domain.Session, error) {

	if duration <= 0 {
		return nil, ErrInvalidDuration
	}

	session := &domain.Session{
		ID:           uuid.New().String(),
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(duration),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Guardar sesión en el repositorio
	err := uc.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// GetSessionByID obtiene una sesión por su ID
func (uc *SessionUseCase) GetSessionByID(ctx context.Context, id string) (*domain.Session, error) {
	session, err := uc.repo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

// GetSessionByToken obtiene una sesión usando el refresh token
func (uc *SessionUseCase) GetSessionByToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	session, err := uc.repo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidSession
	}

	// Verificar si la sesión está bloqueada o expirada
	if session.IsBlocked {
		return nil, ErrSessionBlocked
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// DeleteSession elimina una sesión por ID
func (uc *SessionUseCase) DeleteSession(ctx context.Context, id string) error {
	return uc.repo.DeleteSession(ctx, id)
}

// DeleteSessionsByUserID elimina todas las sesiones de un usuario
func (uc *SessionUseCase) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	return uc.repo.DeleteSessionsByUserID(ctx, userID)
}

// BlockSession bloquea una sesión específica
func (uc *SessionUseCase) BlockSession(ctx context.Context, id string) error {
	session, err := uc.repo.GetSessionByID(ctx, id)
	if err != nil {
		return ErrSessionNotFound
	}

	session.IsBlocked = true
	session.UpdatedAt = time.Now()
	return uc.repo.CreateSession(ctx, session)
}

// UpdateSession actualiza la sesión en el repositorio
func (uc *SessionUseCase) UpdateSession(ctx context.Context, session *domain.Session) error {
	session.UpdatedAt = time.Now() // Actualiza la fecha de modificación
	return uc.repo.UpdateSession(ctx, session)
}
