package repositories

import (
	"context"

	"github.com/google/uuid"
	models "github.com/kevinhc2110/Degree-project-UCP/internal/domain"
)

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	FindByName(ctx context.Context, name string) (*models.Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
}