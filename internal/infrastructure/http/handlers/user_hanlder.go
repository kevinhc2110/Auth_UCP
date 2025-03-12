package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kevinhc2110/Degree-project-UCP/internal/domain"
	"github.com/kevinhc2110/Degree-project-UCP/internal/usecases"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.userUseCase.CreateUser(c.Request.Context(), &user); err != nil {
		if errors.Is(err, usecases.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser maneja la solicitud para actualizar un usuario
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Verificar si el usuario existe antes de actualizarlo
	_, err := h.userUseCase.GetUserByIdentification(c.Request.Context(), user.Identification)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado, no se puede actualizar"})
		return
	}

	if err := h.userUseCase.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}

// func (h *UserHandler) DeleteUser(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
// 		return
// 	}

// 	// Verificar si el usuario existe antes de eliminarlo
// 	_, err = h.userUseCase.GetUserByID(c.Request.Context(), id)
// 	if err != nil {
// 		if errors.Is(err, usecases.ErrUserNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado, no se puede eliminar"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar usuario"})
// 		return
// 	}

// 	if err := h.userUseCase.DeleteUser(c.Request.Context(), id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
// }

// DeleteUser maneja la solicitud para eliminar un usuario
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.userUseCase.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}
