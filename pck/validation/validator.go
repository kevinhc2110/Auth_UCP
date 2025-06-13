package validation

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/kevinhc2110/Auth_UCP/internal/domain"
)

var validate = validator.New()

// Expresiones regulares
var (
	identificationRegex = regexp.MustCompile(`^\d{6,12}$`) // Solo números, entre 6 y 12 dígitos
	nameRegex          = regexp.MustCompile(`^[A-Za-záéíóúÁÉÍÓÚñÑ\s-]{2,50}$`) // Letras y espacios, 2-50 caracteres
)

// ValidatePassword verifica que la contraseña cumpla con los requisitos de seguridad
func ValidatePassword(password string) error {
	var hasUpper, hasDigit, hasSpecial bool

	if len(password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("la contraseña debe incluir al menos una letra mayúscula")
	}
	if !hasDigit {
		return errors.New("la contraseña debe incluir al menos un número")
	}
	if !hasSpecial {
		return errors.New("la contraseña debe incluir al menos un carácter especial")
	}

	return nil
}

// ValidateUser valida los campos de un usuario antes de guardarlo
func ValidateUser(user *domain.User) error {
	// Validar campos obligatorios
	if user.Identification == "" || user.Name == "" || user.Lastname == "" || user.Email == "" || user.Password == "" {
		return errors.New("todos los campos son obligatorios")
	}

	// Validar identificación
	if !identificationRegex.MatchString(user.Identification) {
		return errors.New("la identificación debe contener solo números y tener entre 6 y 12 dígitos")
	}

	// Validar nombre y apellido
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

	// Validar contraseña con la función corregida
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	// Validar que la fecha de último inicio de sesión no sea en el futuro
	if !user.LastLoginAt.IsZero() && user.LastLoginAt.After(time.Now()) {
		return errors.New("la fecha de último inicio de sesión no puede estar en el futuro")
	}

	return nil
}
