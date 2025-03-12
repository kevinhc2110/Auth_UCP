package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
	"github.com/kevinhc2110/Degree-project-UCP/internal/usecases"
)

type sessionRepositorypg struct {
	db *sql.DB
}

// NewSessionRepositorypg crea una nueva instancia del repositorio
func NewSessionRepositorypg(db *sql.DB) *sessionRepositorypg {
	return &sessionRepositorypg{db: db}
}

// CreateSession guarda una nueva sesión en la base de datos
func (r *sessionRepositorypg) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.RefreshToken, session.UserAgent, session.ClientIP,
		session.IsBlocked, session.ExpiresAt, session.CreatedAt, session.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error al crear la sesión: %w", err)
	}
	return nil
}

// GetSessionByID busca una sesión por ID
func (r *sessionRepositorypg) GetSessionByID(ctx context.Context, id string) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at, updated_at FROM sessions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	session := &domain.Session{}
	err := row.Scan(
		&session.ID, &session.UserID, &session.RefreshToken, &session.UserAgent,
		&session.ClientIP, &session.IsBlocked, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, usecases.ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener la sesión: %w", err)
	}
	return session, nil
}

// GetSessionByToken busca una sesión por refresh token
func (r *sessionRepositorypg) GetSessionByToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at, updated_at FROM sessions WHERE refresh_token = $1`
	row := r.db.QueryRowContext(ctx, query, refreshToken)

	session := &domain.Session{}
	err := row.Scan(
		&session.ID, &session.UserID, &session.RefreshToken, &session.UserAgent,
		&session.ClientIP, &session.IsBlocked, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, usecases.ErrInvalidSession
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener sesión por token: %w", err)
	}

	// Verificamos si la sesión está bloqueada o expirada
	if session.IsBlocked {
		return nil, usecases.ErrSessionBlocked
	}
	if session.ExpiresAt.Before(session.CreatedAt) {
		return nil, usecases.ErrSessionExpired
	}

	return session, nil
}

// DeleteSession elimina una sesión por ID
func (r *sessionRepositorypg) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la sesión: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return usecases.ErrSessionNotFound
	}
	return nil
}

// DeleteSessionsByUserID elimina todas las sesiones de un usuario
func (r *sessionRepositorypg) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error al eliminar sesiones del usuario: %w", err)
	}
	return nil
}

// UpdateSession actualiza una sesión existente
func (r *sessionRepositorypg) UpdateSession(ctx context.Context, session *domain.Session) error {
	query := `
		UPDATE sessions
		SET refresh_token = $1, user_agent = $2, client_ip = $3, is_blocked = $4, expires_at = $5, updated_at = $6
		WHERE id = $7
	`
	res, err := r.db.ExecContext(ctx, query,
		session.RefreshToken, session.UserAgent, session.ClientIP,
		session.IsBlocked, session.ExpiresAt, session.UpdatedAt, session.ID,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar la sesión: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return usecases.ErrSessionNotFound
	}
	return nil
}
