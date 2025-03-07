package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// * Configuración SMTP
const (
	SMTPServer = "smtp.gmail.com"
	SMTPPort   = 587
	SMTPUser   = "tuemail@gmail.com"
	SMTPPass   = "tucontraseña" // Usa una clave de aplicación si usas Gmail
)

// * SendRecoveryEmail envía un correo con un enlace de recuperación
func SendRecoveryEmail(toEmail, resetToken string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", SMTPUser)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Recuperación de contraseña")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Hola,</p>
		<p>Haz clic en el siguiente enlace para restablecer tu contraseña:</p>
		<a href="http://tuapp.com/reset-password?token=%s">Restablecer contraseña</a>
	`, resetToken))

	d := gomail.NewDialer(SMTPServer, SMTPPort, SMTPUser, SMTPPass)

	// Enviar correo
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
