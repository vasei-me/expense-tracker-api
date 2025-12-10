package repositories

import (
	"context"
	"expense-tracker/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}