package repository

import "github.com/kevinhc2110/Degree-project-UCP/internal/domain"

type UserRepository interface {
	RegisterUser(user domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	StoreRecoveryToken(userID int64, token string) error
	ChangePassword(email, newPassword string) error
}
