package repositories

import (
	"context"

	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByID(ctx context.Context, id string) (*domain.Session, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	UpdateSession(ctx context.Context, session *domain.Session) error
	DeleteSession(ctx context.Context, id string) error
	DeleteSessionsByUserID(ctx context.Context, userID string) error
}