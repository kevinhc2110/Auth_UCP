package repositories

import (
	"context"

	models "github.com/kevinhc2110/Degree-project-UCP/internal/domain"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteSessionsByUserID(ctx context.Context, userID string) error
}