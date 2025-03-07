package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID        int64     `json:"id"`         // Identificador único
	Name      string    `json:"name"`       // Nombre del usuario
	Lastname  string    `json:"lastname"`   // Apellido del usuario
	Email     string    `json:"email"`      // Correo electrónico
	Password  string    `json:"-"`          // Contraseña (nunca se expone en JSON)
	Role      string    `json:"role"`       // Rol del usuario (admin, user, etc.)
	CreatedAt time.Time `json:"created_at"` // Fecha de creación
	UpdatedAt time.Time `json:"updated_at"` // Fecha de actualización
}

func (u *User) Validate() error {

	// Eliminar espacios en blanco
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	// Validar nombre
	if u.Name == "" {
		return errors.New("el nombre no puede estar vacío")
	}

	// Validar apellido
	if u.Lastname == "" {
		return errors.New("el apellido no puede estar vacío")
	}

	// Validar email con expresión regular
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(u.Email) {
		return errors.New("el email no es válido")
	}

	// Validar contraseña (mínimo 6 caracteres)
	if len(u.Password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	// Validar roles permitidos
	allowedRoles := map[string]bool{"admin": true, "workers": true, "user": true}
	if _, exists := allowedRoles[u.Role]; !exists {
		return errors.New("rol no válido")
	}

	return nil
}

// ValidatePassword verifica que la contraseña sea segura
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\-]`).MatchString(password) // Caracteres especiales comunes

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("la contraseña debe contener al menos una mayúscula, una minúscula, un número y un carácter especial")
	}

	return nil
}

// ValidateEmail verifica si el email tiene un formato válido
func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("el email no es válido")
	}
	return nil
}
