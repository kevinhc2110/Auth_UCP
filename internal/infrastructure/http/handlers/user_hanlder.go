package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kevinhc2110/Auth_UCP/internal/domain"
	"github.com/kevinhc2110/Auth_UCP/internal/usecases"
)

// UserHandler maneja las solicitudes relacionadas con usuarios
type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// CreateUser maneja la solicitud para crear un usuario
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUseCase.CreateUser(c.Request.Context(), &user); err != nil {
		if errors.Is(err, usecases.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente", "user": user})
}

// GetUserByIdentification maneja la solicitud para obtener un usuario por identificación
func (h *UserHandler) GetUserByIdentification(c *gin.Context) {
	identification := c.Param("identification")

	user, err := h.userUseCase.GetUserByIdentification(c.Request.Context(), identification)
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByEmail maneja la solicitud para obtener un usuario por email
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.userUseCase.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser maneja la solicitud para actualizar un usuario
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Obtener el ID del usuario autenticado desde el token JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Validar que el usuario autenticado solo pueda actualizar su propio perfil
	if userID != user.ID.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "No puedes actualizar este usuario"})
		return
	}

	// Verificar si el usuario existe antes de actualizarlo
	existingUser, err := h.userUseCase.GetUserByIdentification(c.Request.Context(), user.Identification)
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado, no se puede actualizar"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Mantener el mismo ID para evitar cambios malintencionados
	user.ID = existingUser.ID

	if err := h.userUseCase.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}


func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Obtener el ID del usuario autenticado desde el token JWT
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// Convertir userID a uuid.UUID
	userID, err := uuid.Parse(fmt.Sprintf("%v", userIDRaw))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
		return
	}

	// Obtener el rol del usuario autenticado
	role, _ := c.Get("role")

	// Obtener el ID del usuario a eliminar desde la URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar si el usuario autenticado es el mismo que intenta eliminar o si es admin
	if userID != id && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No puedes eliminar este usuario"})
		return
	}

	// Verificar si el usuario existe antes de eliminarlo
	_, err = h.userUseCase.GetUserByID(c.Request.Context(), id) 
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Eliminar usuario
	if err := h.userUseCase.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}
