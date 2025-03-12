package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
	"github.com/kevinhc2110/Degree-project-UCP/internal/repositories"
	"github.com/kevinhc2110/Degree-project-UCP/internal/usecases"
	"github.com/lib/pq"
)

type UserRepositoryPg struct {
	db *sql.DB
}

func NewUserRepositoryPg(db *sql.DB) repositories.UserRepository {
	return &UserRepositoryPg{db: db}
}

func (r *UserRepositoryPg) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, identification, name, lastname, email, password, role, active, created_at, updated_at, lastlogin_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Identification, user.Name, user.Lastname, user.Email, user.Password, user.Role,
		user.Active, user.CreatedAt, user.UpdatedAt, user.LastLoginAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" { // Código de error para violación de restricción UNIQUE
				return usecases.ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("error al crear el usuario en la base de datos: %w", err)
	}
	return nil
}

func (r *UserRepositoryPg) FindByIdentification(ctx context.Context, identification string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, identification, email, password, active, created_at, updated_at FROM users WHERE identification = $1`
	err := r.db.QueryRowContext(ctx, query, identification).Scan(
		&user.ID, &user.Identification, &user.Email, &user.Password,
		&user.Active, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, usecases.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario por identificación: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryPg) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, identification, email, password, active, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Identification, &user.Email, &user.Password,
		&user.Active, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, usecases.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario por email: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryPg) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET identification = $1, email = $2, password = $3, active = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, user.Identification, user.Email, user.Password, user.Active, user.UpdatedAt, user.ID)

	if err != nil {
		return errors.New("error al actualizar el usuario")
	}
	return nil
}

func (r *UserRepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return errors.New("error al eliminar el usuario")
	}
	return nil
}
