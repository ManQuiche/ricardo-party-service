package app

import (
	"context"
	"ricardo/party-service/internal/core/entities"
	"ricardo/party-service/internal/core/ports"
)

type UserService interface {
	ports.UserService
}

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) UserService {
	return userService{
		repo: repo,
	}
}

func (p userService) Get(ctx context.Context, userID uint) (*entities.User, error) {
	return p.repo.Get(ctx, userID)
}

func (p userService) GetAll(ctx context.Context) ([]entities.User, error) {
	return p.repo.GetAll(ctx)
}

func (p userService) Save(ctx context.Context, user entities.User) (*entities.User, error) {
	return p.repo.Save(ctx, user)
}

func (p userService) Delete(ctx context.Context, userID uint) error {
	return p.repo.Delete(ctx, userID)
}
