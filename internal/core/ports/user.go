package ports

import (
	"context"
	"ricardo/party-service/internal/core/entities"
)

type UserService interface {
	Get(ctx context.Context, userID uint) (*entities.User, error)
	GetAll(ctx context.Context) ([]entities.User, error)
	Save(ctx context.Context, user entities.User) (*entities.User, error)
	Delete(ctx context.Context, userID uint) error
}

type UserRepository interface {
	Get(ctx context.Context, userID uint) (*entities.User, error)
	GetAll(ctx context.Context) ([]entities.User, error)
	Save(ctx context.Context, user entities.User) (*entities.User, error)
	Delete(ctx context.Context, userID uint) error
}
