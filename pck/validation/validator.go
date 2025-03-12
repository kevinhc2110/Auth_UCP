package validation

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
)

var validate = validator.New()

// Expresiones regulares
var (
	passwordRegex      = regexp.MustCompile(`^(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
	identificationRegex = regexp.MustCompile(`^\d{6,12}$`) // Solo números, entre 6 y 12 dígitos
	nameRegex          = regexp.MustCompile(`^[A-Za-záéíóúÁÉÍÓÚñÑ\s-]{2,50}$`) // Letras y espacios, 2-50 caracteres
)

// ValidateUser valida los campos de un usuario antes de guardarlo
func ValidateUser(user *domain.User) error {
	// Validar campos obligatorios
	if user.Identification == "" || user.Name == "" || user.Lastname == "" || user.Email == "" || user.Password == "" {
		return errors.New("todos los campos son obligatorios")
	}

	// Validar identificación (solo números, entre 6 y 12 dígitos)
	if !identificationRegex.MatchString(user.Identification) {
		return errors.New("la identificación debe contener solo números y tener entre 6 y 12 dígitos")
	}

	// Validar nombre y apellido (solo letras, espacios y guiones)
	if !nameRegex.MatchString(user.Name) {
		return errors.New("el nombre solo puede contener letras y espacios, con un mínimo de 2 caracteres")
	}
	if !nameRegex.MatchString(user.Lastname) {
		return errors.New("el apellido solo puede contener letras y espacios, con un mínimo de 2 caracteres")
	}

	// Validar correo electrónico
	if err := validate.Var(user.Email, "required,email"); err != nil {
		return errors.New("el correo electrónico no es válido")
	}

	// Validar contraseña segura
	if !passwordRegex.MatchString(user.Password) {
		return errors.New("la contraseña debe tener al menos 8 caracteres, una mayúscula, un número y un carácter especial")
	}

	// Validar que la fecha de último inicio de sesión no sea en el futuro
	if !user.LastLoginAt.IsZero() && user.LastLoginAt.After(time.Now()) {
		return errors.New("la fecha de último inicio de sesión no puede estar en el futuro")
	}

	return nil
}
