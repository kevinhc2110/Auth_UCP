package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
	"github.com/kevinhc2110/Degree-project-UCP/pck/validation"

	"github.com/kevinhc2110/Degree-project-UCP/internal/infrastructure/security"
	"github.com/kevinhc2110/Degree-project-UCP/internal/repositories"
)

var (
	ErrUserNotFound       = errors.New("usario no encontrado")
	ErrEmailAlreadyExists = errors.New("el email ya esta registrado en una cuenta")
	ErrInvalidCredentials = errors.New("credenciales invalidas")
)

type UserUseCase struct {
	repo repositories.UserRepository
}

// NewUserUseCase crea una nueva instancia de UserUseCase
func NewUserUseCase(repo repositories.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *domain.User) error {
	// Validar usuario
	if err := validation.ValidateUser(user); err != nil {
		return err
	}

	// Verificar si el correo ya está registrado
	if existingUser, _ := uc.repo.FindByEmail(ctx, user.Email); existingUser != nil {
		return ErrEmailAlreadyExists
	}

	// Cifrar la contraseña
	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return errors.New("error al cifrar la contraseña")
	}
	user.Password = hashedPassword

	// Asignar valores predeterminados
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Active = true

	// Guardar el usuario en la base de datos
	if err := uc.repo.Create(ctx, user); err != nil {
		return errors.New("error al guardar el usuario")
	}

	return nil
}

// GetUserByID obtiene un usuario por su ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	user, err := uc.repo.FindByID(ctx, ID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByIdentification obtiene un usuario por su identification
func (uc *UserUseCase) GetUserByIdentification(ctx context.Context, identification string) (*domain.User, error) {
	user, err := uc.repo.FindByIdentification(ctx, identification)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByEmail obtiene un usuario por su correo electrónico
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser actualiza la información de un usuario
func (uc *UserUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, user); err != nil {
		return errors.New("error al actualizar el usuario")
	}
	return nil
}

// DeleteUser elimina un usuario por su ID
func (uc *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return errors.New("error al eliminar el usuario")
	}
	return nil
}
