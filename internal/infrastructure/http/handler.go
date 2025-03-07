package http

import (
	"encoding/json"
	"net/http"

	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
	"github.com/kevinhc2110/Degree-project-UCP/internal/usecase"
)

// AuthHandler maneja las solicitudes relacionadas con autenticación
type AuthHandler struct {
	authUC *usecase.AuthUseCase
}

// NewAuthHandler crea un nuevo AuthHandler
func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// RegisterHandler maneja el registro de usuarios
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := h.authUC.Register(&user); err != nil {
		http.Error(w, "Error al registrar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado correctamente"})
}

// GetUserByEmailHandler maneja la solicitud para obtener un usuario por su email
func (h *AuthHandler) GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Autenticar al usuario y obtener el token
	token, err := h.authUC.Authenticate(request.Email, request.Password)
	if err != nil {
		http.Error(w, "Credenciales inválidas: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Devolver el token en la respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RememberPasswordHandler maneja la solicitud para recuperar contraseña
func (h *AuthHandler) RememberPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := h.authUC.GenerateRecoveryToken(request.Email); err != nil {
		http.Error(w, "Error al recuperar la contraseña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Correo de recuperación enviado"})
}

// ChangePasswordHandler maneja el cambio de contraseña
func (h *AuthHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := h.authUC.ChangePassword(request.Email, request.NewPassword); err != nil {
		http.Error(w, "Error al cambiar la contraseña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Contraseña cambiada correctamente"})
}
