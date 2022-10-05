package postgresql

import (
	"context"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"gitlab.com/ricardo134/party-service/internal/core/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB) ports.UserRepository {
	return userRepository{client: client}
}

func (p userRepository) Get(ctx context.Context, userID uint) (*entities.User, error) {
	var user entities.User
	err := p.client.First(&user, userID).Error
	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return &user, nil
}

func (p userRepository) GetAll(ctx context.Context) ([]entities.User, error) {
	var parties []entities.User
	err := p.client.Find(&parties).Error
	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return parties, nil
}

func (p userRepository) Save(ctx context.Context, user entities.User) (*entities.User, error) {
	err := p.client.Save(&user).Error

	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return &user, err
}

func (p userRepository) Delete(ctx context.Context, userID uint) error {
	err := p.client.Delete(&entities.User{}, userID).Error
	if err != nil {
		return notFoundOrElseError(err)
	}

	return nil
}
