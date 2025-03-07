package usecase

import (
	"errors"
	"time"

	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
	"github.com/kevinhc2110/Degree-project-UCP/internal/repository"
	"github.com/kevinhc2110/Degree-project-UCP/pkg/hash"
	"github.com/kevinhc2110/Degree-project-UCP/pkg/mail"
	"github.com/kevinhc2110/Degree-project-UCP/pkg/token"
)

// AuthUseCase maneja la lógica de autenticación y usuarios
type AuthUseCase struct {
	userRepo repository.UserRepository
}

// NewAuthUseCase crea una nueva instancia de AuthUseCase
func NewAuthUseCase(userRepo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

// Register registra un nuevo usuario
func (uc *AuthUseCase) Register(user *domain.User) error {

	// Asignar el rol "user" por defecto
	user.Role = "user"

	// Asignar fechas de creación y actualización
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Validar usuario antes de registrarlo
	if err := user.Validate(); err != nil {
		return err
	}

	// Hashear la contraseña antes de guardarla
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return errors.New("error al encriptar la contraseña")
	}
	user.Password = hashedPassword

	// Registrar usuario en el repositorio
	return uc.userRepo.RegisterUser(*user)
}

// Authenticate valida credenciales y devuelve un token JWT si son correctas
func (uc *AuthUseCase) Authenticate(email, password string) (string, error) {

	// Validar formato de email antes de consultar la base de datos
	if err := domain.ValidateEmail(email); err != nil {
		return "", err
	}

	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("usuario no encontrado")
	}

	// Comparar contraseñas
	if !hash.CheckPasswordHash(password, user.Password) {
		return "", errors.New("contraseña incorrecta")
	}

	// Generar token JWT
	tokenString, err := token.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", errors.New("error al generar el token")
	}

	return tokenString, nil
}

// RememberPassword envía un email con instrucciones para recuperar la contraseña
func (uc *AuthUseCase) GenerateRecoveryToken(email string) error {
	// Verificar si el usuario existe en la base de datos
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("correo no registrado")
	}

	// Generar un token de recuperación válido por 1 hora
	resetToken, err := token.GenerateTokenMail(user.ID, user.Email, "reset")
	if err != nil {
		return errors.New("error al generar el token de recuperación")
	}

	// Guardar el token en la base de datos
	err = uc.userRepo.StoreRecoveryToken(user.ID, resetToken)
	if err != nil {
		return errors.New("error al almacenar el token de recuperación")
	}

	// Enviar el email con el token
	if err := mail.SendRecoveryEmail(email, resetToken); err != nil {
		return errors.New("error al enviar el correo de recuperación")
	}

	return nil
}

// ChangePassword cambia la contraseña de un usuario autenticado
func (uc *AuthUseCase) ChangePassword(email, newPassword string) error {
	// Validar que la nueva contraseña
	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hashear la nueva contraseña
	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return errors.New("error al encriptar la nueva contraseña")
	}

	// Actualizar contraseña en el repositorio
	return uc.userRepo.ChangePassword(email, hashedPassword)
}

// ResetPassword permite a un usuario cambiar su contraseña usando un token de recuperación
func (uc *AuthUseCase) ResetPassword(resetToken, newPassword string) error {
	// Validar que la contraseña cumpla con los requisitos
	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Verificar la validez del token y extraer los datos del usuario
	claims, err := token.ValidateToken(resetToken)
	if err != nil {
		return errors.New("token inválido o expirado")
	}

	// Hashear la nueva contraseña
	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return errors.New("error al encriptar la nueva contraseña")
	}

	// Actualizar la contraseña en la base de datos
	if err := uc.userRepo.ChangePassword(claims.Email, hashedPassword); err != nil {
		return errors.New("error al actualizar la contraseña")
	}

	return nil
}
