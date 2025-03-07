package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
)

// UserRepositoryDB es la implementación de UserRepository con base de datos
type UserRepositoryDB struct {
	db *sql.DB
}

// NewUserRepositoryDB crea una nueva instancia de UserRepositoryDB
func NewUserRepositoryDB(db *sql.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}

// RegisterUser guarda un nuevo usuario en la base de datos
func (r *UserRepositoryDB) RegisterUser(user domain.User) error {
	query := `INSERT INTO users (name, lastname, email, password, role, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := r.db.QueryRow(query, user.Name, user.Lastname, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).
		Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail busca un usuario por su email
func (r *UserRepositoryDB) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, lastname, email, password, role, created_at, updated_at FROM users WHERE email = ?`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Lastname, &user.Email,
		&user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

// StoreRecoveryToken guarda un token de recuperación en la base de datos
func (r *UserRepositoryDB) StoreRecoveryToken(userID int64, token string) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE 
		SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at
	`

	expirationTime := time.Now().Add(1 * time.Hour) // Token válido por 1 hora
	_, err := r.db.Exec(query, userID, token, expirationTime)
	return err
}

// ChangePassword actualiza la contraseña de un usuario
func (r *UserRepositoryDB) ChangePassword(email, newPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE email = ?`

	updatedAt := time.Now()
	result, err := r.db.Exec(query, newPassword, updatedAt, email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("usuario no encontrado o contraseña no cambiada")
	}

	return nil
}
